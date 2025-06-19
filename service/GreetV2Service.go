package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/global/enum"
	"git.exclouds.org/ww/greet_ws/model"
	"git.exclouds.org/ww/greet_ws/model/req"
	"git.exclouds.org/ww/greet_ws/model/resp"
	"git.exclouds.org/ww/greet_ws/utils"
	"github.com/tidwall/gjson"
)

type GreetService struct {
	mu sync.Mutex
}

func NewGreetService() *GreetService {
	return &GreetService{}
}

func (svc *GreetService) GreetTaskV2() {

	var httpClientUtils = utils.HttpClientUtils{
		Header: map[string][]string{"token": {}},
	}

	//获取租户列表及配置信息
	respWrapper := httpClientUtils.Get(GetTenantList, HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		global.LOGGER.Error(fmt.Sprintf("获取租户列表失败:%s, %s \n", GetTenantList, respWrapper.Body))
		return
	}

	respBody := respWrapper.Body
	msgCode := gjson.Get(respBody, "code").Int()
	if msgCode != 200 {
		global.LOGGER.Error(fmt.Sprintf("获取租户账号: %s\n", respBody))
		return
	}

	var tenantDataList []resp.TenantModel

	global.LOGGER.Info("tenant.respBody:" + respBody)
	tenantJson := gjson.Get(respBody, "data").String()
	if tenantJson == "" {
		global.LOGGER.Info(fmt.Sprintf("获取租户账号为空: %s\n", respBody))
		time.Sleep(time.Second)
		return
	}

	if err := json.Unmarshal([]byte(tenantJson), &tenantDataList); err != nil {
		global.LOGGER.Error(fmt.Sprintf("租户格式不正确：%s\n", tenantJson))
		return
	}

	for _, tenant := range tenantDataList {
		go svc.scheduleTask(tenant)
	}
}

// 模拟定时任务执行
func (svc *GreetService) scheduleTask(tenant resp.TenantModel) {

	var httpClientUtils = utils.HttpClientUtils{
		Header: map[string][]string{"token": {svc.GlobalTaskConfig.Token}},
	}

	//获取联系人账号
	reqJson := fmt.Sprintf(`{"userId": "%s","count": %d}`, tenant.ID, tenant.MessagesPerMinute)

	respWrapper := httpClientUtils.PutJson(GetTargetList, reqJson, HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		global.LOGGER.Error(fmt.Sprintf("获取联系人列表失败:%s, %s \n", GetTargetList, respWrapper.Body))
	}

	respBody := respWrapper.Body
	msgCode := gjson.Get(respBody, "code").Int()
	if msgCode != 200 {
		global.LOGGER.Error(fmt.Sprintf("获取联系人列表: %s\n", respBody))
	}

	var targetList []resp.TargetModel

	global.LOGGER.Info("target.respBody:" + respBody)
	targetJson := gjson.Get(respBody, "data").String()
	if targetJson == "" {
		global.LOGGER.Info(fmt.Sprintf("获取联系人为空: %s\n", respBody))
	}

	if err := json.Unmarshal([]byte(targetJson), &targetList); err != nil {
		global.LOGGER.Error(fmt.Sprintf("联系人格式不正确：%s\n", targetJson))
	}

	// 判断切片是否为空
	if len(targetList) != 0 {
		for _, target := range targetList {
			go svc.sendMess(tenant, target)
		}
	}
	global.LOGGER.Info(fmt.Sprintf("tenantID：%s 下联系人列表为空,本次任务结束", tenant.ID))
}

func (svc *GreetService) sendMess(tenant resp.TenantModel, target resp.TargetModel) {
	// 获取初始账号数据
	var account *model.PullAccountModelV2

	// 记录开始时间
	startTime := time.Now()
	for {
		// 检查是否超过十分钟
		if time.Since(startTime) > 10*time.Minute {
			global.LOGGER.Info(fmt.Sprintf("租户id：%s ===> 坐席id：%s ===> 目标 %s：发送消息超时，已超过 10 分钟，抛弃", tenant.ID, target.BelongToSubUserID, target.Phone))
			break
		}
		// 如果账号为空，重新拉取账号
		if account == nil {
			account, _ = svc.GetWsAccount(target.BelongToSubUserID, target.CountryCode)
			if account == nil {
				global.LOGGER.Info(fmt.Sprintf("租户id：%s ===> 坐席id：%s ===> 目标 %s：没有可用账号，20 秒后重试", tenant.ID, target.BelongToSubUserID, target.Phone))
				time.Sleep(20 * time.Second)
				continue
			}
			global.LOGGER.Info(fmt.Sprintf("租户id：%s ===> 坐席id：%s ===> 目标 %s,拉取到新ws账号: %s。", tenant.ID, target.BelongToSubUserID, target.Phone, account.Username))
		}
		// 使用可用账号发送消息
		err := svc.sendMessage(tenant, target, account)
		// 更新账号的可用次数
		if err == nil {
			break
		} else {
			global.LOGGER.Error(fmt.Sprintf("发送消息失败: %v,重新拉取账号", err))
			account = nil
			continue
		}
	}

}

func (svc *GreetService) finishAccount(tenantID string, seatId string, wsAccountId string, messageId string, wsAccount string, targetId string, targetPhone string, failReason string) (err error) {
	var httpClientUtils = utils.HttpClientUtils{
		Header: map[string][]string{"token": {svc.GlobalTaskConfig.Token}},
	}

	// 修复 reqBody 中的 bool 变量格式化
	reqBody := fmt.Sprintf(`{"parentId":"%s","userId": "%s","wsAccountId":"%s","messageId":"%s","wsAccount": "%s","targetId":"%s","targetPhone":"%s","failReason":"%s"}`, tenantID, seatId, wsAccountId, messageId, wsAccount, targetId, targetPhone, failReason)
	respWrapper := httpClientUtils.PutJson(WsTaskFinish, reqBody, HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("上报ws账号任务完成失败:%s, %s", WsTaskFinish, respWrapper.Body)
		global.LOGGER.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	respBody := respWrapper.Body
	msgCode := gjson.Get(respBody, "code").Int()
	if msgCode != 200 {
		errMsg := fmt.Sprintf("上报ws账号任务完成失败: %s", respBody)
		global.LOGGER.Error(errMsg)
		return fmt.Errorf(errMsg)
	}
	global.LOGGER.Info(fmt.Sprintf("上报ws账号任务完成:%+v", reqBody))
	return nil
}

// sendMessage 发送消息
func (svc *GreetService) sendMessage(tenant resp.TenantModel, target resp.TargetModel, account *model.PullAccountModelV2) (err error) {
	taskContent := tenant.GreetMessageContent

	// 使用 strconv.ParseInt 直接将字符串转换为 int64，并处理错误
	seatId, err := strconv.ParseInt(target.BelongToSubUserID, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("坐席id转换错误:%s", err)
		global.LOGGER.Error(errMsg)
		return err
	}

	wsMsgModel := WsMsgModelV2{
		UUID:        account.Username,
		ToUserId:    target.Phone,
		TaskContent: taskContent,
		SeatID:      seatId,
		Type:        account.Type,
	}

	// 调用发送方法
	err = svc.greetV2(tenant.ID, target.ID, wsMsgModel, account.Id, account.NodeId)
	return
}

// 随机取一条方法
func getRandomGreet(greetTextList []string) string {
	// 检查列表是否为空
	if len(greetTextList) == 0 {
		return "No greetings available."
	}

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 随机索引
	randomIndex := rand.Intn(len(greetTextList))
	return greetTextList[randomIndex]
}

func (svc *GreetService) GetWsAccount(belongToSubUserId string, countryCode string) (*model.PullAccountModelV2, error) {
	var httpClientUtils = utils.HttpClientUtils{
		Header: map[string][]string{"token": {svc.GlobalTaskConfig.Token}},
	}

	// 获取 ws 账号列表
	reqBody := fmt.Sprintf(`{"userId":"%s", "countryCode":"%s"}`, belongToSubUserId, countryCode)
	respWrapper := httpClientUtils.PutJson(GetWsAccountList, reqBody, HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("获取账号列表失败:%s, %s", GetWsAccountList, respWrapper.Body)
		global.LOGGER.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	respBody := respWrapper.Body
	msgCode := gjson.Get(respBody, "code").Int()
	if msgCode != 200 {
		errMsg := fmt.Sprintf("获取号池账号失败: %s", respBody)
		global.LOGGER.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	var pullAccount model.PullAccountModelV2
	pullAccountData := gjson.Get(respBody, "data").Raw
	if err := json.Unmarshal([]byte(pullAccountData), &pullAccount); err != nil {
		errMsg := fmt.Sprintf("解析号池账号 JSON 失败: %s", respBody)
		global.LOGGER.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	if len(pullAccount.Id) == 0 || len(pullAccount.Username) == 0 {
		return nil, fmt.Errorf("获取账号为空:%s, %s", GetWsAccountList, respBody)
	}
	// 返回获取的账号
	return &pullAccount, nil
}

type WsMsgModelV2 struct {
	UUID        string //uuid 是 ws账号。
	ToUserId    string //touserid 是联系人手机号
	TaskContent *resp.GreetMessageContent
	SeatID      int64
	Type        int
}

// greet 打招呼
func (svc *GreetService) greetV2(tenantID string, targetId string, wsMsgModel WsMsgModelV2, wsAccountId string, nodeId string) (err error) {
	seatId := wsMsgModel.SeatID                              //下发 座席id 分表等有用
	handledUuid := strings.TrimPrefix(wsMsgModel.UUID, "QR") //对QR的ws账号进行兼容处理
	waAccount := handledUuid
	UUID := handledUuid
	toUserId := wsMsgModel.ToUserId
	phone := wsMsgModel.ToUserId
	waAccountType := wsMsgModel.Type

	wsAccountIdIntValue, err := strconv.ParseInt(wsAccountId, 10, 64)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("转换失败:", err))
	}

	contactGroup := phone + waAccount
	rds := global.REDIS
	rdsKey := "contact_group:" + contactGroup
	exists, err := redisUtils.KeyExists(rds, rdsKey)
	if exists {
		errMsg := fmt.Sprintf("目标:%s + ws账号:%s 组合已存在，ws账号跳过并上报失败处理", phone, waAccount)
		global.LOGGER.Error(errMsg)
		svc.finishAccount(tenantID, strconv.FormatInt(seatId, 10), wsAccountId, "", waAccount, "", "", "")
		return errors.New(errMsg)
	}

	wsAddMessageModel, err := svc.Initsession(wsMsgModel.UUID, phone)
	if err != nil {
		global.LOGGER.Error("初始化会话失败：" + err.Error())
		_ = svc.finishAccount(tenantID, strconv.FormatInt(seatId, 10), wsAccountId, "", waAccount, "", "", wsAddMessageModel.CgiBaseResponse.ErrMsg)
		return
	}
	global.LOGGER.Info("初始化会话成功，结果：" + fmt.Sprintf("%+v", wsAddMessageModel))
	contactID := global.GenerateUniqueId()
	contactModel := model.ContactModel{
		ID:              contactID,
		DisplayName:     phone,
		Avatar:          "",
		Index:           "0",
		Unread:          0,
		Phone:           phone,
		ContactID:       contactID,
		WaAccountID:     wsAccountIdIntValue,
		WaAccount:       waAccount,
		WaAccountTo:     toUserId,
		AllowReply:      0,
		IsOnline:        1,
		IsCustomerBlock: 0,
		IsPinTop:        0,
		//ClientVer:       1,
		//CanUseCallOffer: 0,
		LastSendTime: time.Now().Unix(),
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
		GroupID:      "",
		UUID:         UUID,
		ToUserId:     toUserId,
		SeatID:       seatId,
		NodeID:       nodeId,
	}

	messageType := wsMsgModel.TaskContent.Type //enum.MessageTypeText

	cgiRequest := req.CgiRequest{
		ToUserId:  toUserId,
		IsGroup:   false,
		IsFakeMsg: false,
	}

	if messageType == "text" {
		textList := wsMsgModel.TaskContent.TextList
		reqModel := req.WsSendMessageTextModel{
			WithCgiRequest: req.WithCgiRequest{
				CgiRequest:  cgiRequest,
				StanzaId:    "",
				Participant: "",
				Content:     getRandomGreet(textList),
			},
		}
		//ws发消息，打招呼(文本)
		wsSendMessageModel, err := svc.WsSendAndSaveMessageV2(reqModel, contactModel, req.SendMessageModel{Type: messageType}, tenantID, targetId, wsAccountId, waAccountType)
		if err != nil {
			global.LOGGER.Error("文本===>打招呼失败：" + err.Error())
			_ = svc.finishAccount(tenantID, strconv.FormatInt(seatId, 10), wsAccountId, "", waAccount, "", "", wsSendMessageModel.CgiBaseResponse.ErrMsg)
			return err
		}
		global.LOGGER.Info("文本===>打招呼成功，结果：" + fmt.Sprintf("%+v", wsSendMessageModel))
	} else if messageType == "bot" {
		sourceUrlList := wsMsgModel.TaskContent.UrlList
		reqModel := req.WsSendMessageRichTextModel{
			req.RichTextCgiRequest{
				ToUserId:  toUserId,
				IsGroup:   false,
				Title:     wsMsgModel.TaskContent.Title,
				Text:      wsMsgModel.TaskContent.SecondaryTitle,
				Body:      wsMsgModel.TaskContent.Content,
				ImageUrl:  wsMsgModel.TaskContent.ImageUrl,
				SourceUrl: getRandomGreet(sourceUrlList),
			},
		}
		wsSendMessageModel, err := svc.WsSendAndSaveRichTextMessageV2(reqModel, contactModel, req.SendMessageModel{Type: messageType}, tenantID, targetId, wsAccountId)
		if err != nil {
			global.LOGGER.Error("图文===>打招呼失败：" + err.Error())
			_ = svc.finishAccount(tenantID, strconv.FormatInt(seatId, 10), wsAccountId, "", waAccount, "", "", wsSendMessageModel.CgiBaseResponse.ErrMsg)
			return err
		}
		global.LOGGER.Info("图文===>打招呼成功，结果：" + fmt.Sprintf("%+v", wsSendMessageModel))
	}

	if waAccountType == 2 || waAccountType == 3 {
		go svc.HandleHijack(rds, UUID, waAccountType, phone)
	}

	//保存到会话
	contactService.SeatID = seatId
	err = contactService.Add(contactModel)
	if err != nil {
		global.LOGGER.Error("保存到会话表失败：" + err.Error())
		return
	}
	redisUtils.Set(rds, rdsKey, seatId, time.Hour*24*8)
	global.LOGGER.Info(fmt.Sprintf("contact_group缓存设置成功,key:%s,val:%s", rdsKey, seatId))
	return nil
}

// 调ws发消息，并保存(文本)
func (svc *GreetService) WsSendAndSaveMessageV2(reqModel interface{}, contactModel model.ContactModel, reqSendMessageModel req.SendMessageModel, tenantID string, targetId string, wsAccountId string, accType int) (resp.WsMessageModel, error) {
	var err error
	var wsMessageModel resp.WsMessageModel
	uuid := contactModel.UUID
	if accType == 3 {
		uuid = "QR" + contactModel.UUID
	}

	messageType := reqSendMessageModel.Type
	sendUrl := ""
	content := ""
	switch messageType {
	case enum.MessageTypeText: //文字消息
		tmp := reqModel.(req.WsSendMessageTextModel)
		content = tmp.WithCgiRequest.Content
		sendUrl = WsSendTxtUrl

	case enum.MessageTypeMedia: //图片消息
		tmp := reqModel.(req.WsSendMessageImageUrlModel)
		content = tmp.ImageCgiRequest.Content //"图片文字"
		sendUrl = WsSendImageUrl

	default:
		return wsMessageModel, errors.New("消息类型不支持：" + messageType)
	}

	//组装请求ws数据
	bt, err := json.Marshal(reqModel)
	if err != nil {
		return wsMessageModel, err
	}

	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	headers := map[string]string{
		"RegistrationID": contactModel.NodeID,
	}
	global.LOGGER.Info(fmt.Sprintf("消息请求路径:%s,消息请求header：%+v,消息请求体:%s", endpointUrl, headers, string(bt)))
	respWrapper := httpClientUtils.PostJsonWithHeaders(endpointUrl, string(bt), headers, HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New(respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		return wsMessageModel, fmt.Errorf("请求ws失败：%s", wsMessageModel.CgiBaseResponse.ErrMsg)
	}

	messageId := ""
	// 类型断言，确保 ResponseData 是 map[string]interface{}
	if responseData, ok := wsMessageModel.ResponseData.(map[string]interface{}); ok {
		// 从 map 中获取 "id" 字段
		messageId = responseData["ID"].(string)
		global.LOGGER.Info(fmt.Sprintf("消息ID:", messageId))
	}

	// 在调用 finishWsAccount 的地方需要传入 bool 参数
	_ = svc.finishAccount(tenantID, strconv.FormatInt(contactModel.SeatID, 10), wsAccountId, messageId, contactModel.WaAccount, targetId, contactModel.WaAccountTo, "")

	messageModel := model.MessageModel{
		Ack:                  1,
		ID:                   messageId,
		To:                   contactModel.ID,
		Status:               "0",
		ToContactID:          contactModel.ID,
		SendTime:             time.Now().Unix(),
		FromUser_ID:          "1",
		FromUser_DisplayName: "",
		FromUser_Avatar:      "",
		IsFromMe:             1,
		FromLang:             "",
		ToLang:               "",
		Type:                 messageType,
		AesKey:               "",
		Content:              content,
		ContentTranslate:     "",
		SeatID:               contactModel.SeatID,
		CreateTime:           time.Now().Unix(),
		UpdateTime:           time.Now().Unix(),
		MessageID:            messageId,
		FromJID:              contactModel.WaAccountTo,
		UserJID:              contactModel.WaAccount,
	}
	messageService.SeatID = contactModel.SeatID
	err = messageService.Add(messageModel)
	if err != nil {
		global.LOGGER.Error("保存到消息列表失败：" + err.Error())
		return wsMessageModel, err
	}
	global.LOGGER.Info(fmt.Sprintf("消息保存成功:%+v", messageModel))
	return wsMessageModel, nil
}

/*
*富文本
 */
func (svc *GreetService) WsSendAndSaveRichTextMessageV2(reqModel req.WsSendMessageRichTextModel, contactModel model.ContactModel, reqSendMessageModel req.SendMessageModel, tenantID string, targetId string, wsAccountId string) (resp.WsMessageModel, error) {
	var err error
	var wsMessageModel resp.WsMessageModel
	messageType := enum.MessageTypeRichText
	sendUrl := WsSendRichTextUrl
	//组装请求ws数据
	bt, err := json.Marshal(reqModel)
	if err != nil {
		global.LOGGER.Error("reqModel 转换失败：" + err.Error())
		return wsMessageModel, err
	}

	contentModel := model.ContentRichTextModel{
		Text: reqModel.WithCgiRequest.Text,
		//Type: enum.RichTextTypeHyperlink,
		//HyperlinkContextInfo: struct {
		//	Title     string `json:"title"`
		//	Body      string `json:"body"`
		//	ImageUrl  string `json:"imageUrl"`
		//	SourceUrl string `json:"sourceUrl"`
		//}{
		//	Title:     reqModel.WithCgiRequest.Title,
		//	Body:      reqModel.WithCgiRequest.Body,
		//	ImageUrl:  reqModel.WithCgiRequest.ImageUrl,
		//	SourceUrl: reqModel.WithCgiRequest.SourceUrl,
		//},
	}
	contentBt, err := json.Marshal(contentModel)
	if err != nil {
		global.LOGGER.Error("打招呼，图文contentModel 转换失败：" + err.Error())
		return wsMessageModel, err
	}
	content := string(contentBt)
	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, contactModel.UUID)
	headers := map[string]string{
		"RegistrationID": contactModel.NodeID,
	}
	global.LOGGER.Info(fmt.Sprintf("消息请求路径:%s,消息请求header：%+v,消息请求体:%s", endpointUrl, headers, string(bt)))
	respWrapper := httpClientUtils.PostJsonWithHeaders(endpointUrl, string(bt), headers, HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New(respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		return wsMessageModel, fmt.Errorf("请求ws失败：%s", wsMessageModel.CgiBaseResponse.ErrMsg)
	}

	messageId := ""
	// 类型断言，确保 ResponseData 是 map[string]interface{}
	if responseData, ok := wsMessageModel.ResponseData.(map[string]interface{}); ok {
		// 从 map 中获取 "id" 字段
		messageId = responseData["ID"].(string)
		global.LOGGER.Info(fmt.Sprintf("消息ID:", messageId))
	}

	// 在调用 finishWsAccount 的地方需要传入 bool 参数
	_ = svc.finishAccount(tenantID, strconv.FormatInt(contactModel.SeatID, 10), wsAccountId, messageId, contactModel.WaAccount, targetId, contactModel.WaAccountTo, "")

	messageModel := model.MessageModel{
		Ack:                  1,
		ID:                   messageId,
		To:                   contactModel.ID,
		Status:               "0",
		ToContactID:          contactModel.ID,
		SendTime:             time.Now().Unix(),
		FromUser_ID:          "1",
		FromUser_DisplayName: "",
		FromUser_Avatar:      "",
		IsFromMe:             1,
		FromLang:             "",
		ToLang:               "",
		Type:                 messageType,
		AesKey:               "",
		Content:              content,
		ContentTranslate:     "",
		SeatID:               contactModel.SeatID,
		CreateTime:           time.Now().Unix(),
		UpdateTime:           time.Now().Unix(),
		MessageID:            messageId,
		FromJID:              contactModel.WaAccountTo,
		UserJID:              contactModel.WaAccount,
	}
	messageService.SeatID = contactModel.SeatID
	err = messageService.Add(messageModel)
	if err != nil {
		global.LOGGER.Error("保存到消息列表失败：" + err.Error())
		return wsMessageModel, err
	}
	global.LOGGER.Info(fmt.Sprintf("消息保存成功:%+v", messageModel))
	return wsMessageModel, nil
}
func (svc *GreetService) WsSendWsMessageOnly(reqModel interface{}, contactModel model.ContactModel, reqSendMessageModel req.SendMessageModel) (resp.WsMessageModel, error) {
	var err error
	var wsMessageModel resp.WsMessageModel
	messageType := reqSendMessageModel.Type
	sendUrl := ""
	switch messageType {
	case enum.MessageTypeText: //文字消息
		sendUrl = WsSendTxtUrl

	case enum.MessageTypeMedia: //图片消息
		sendUrl = WsSendImageUrl

	default:
		return wsMessageModel, errors.New("消息类型不支持：" + messageType)
	}

	//组装请求ws数据
	bt, err := json.Marshal(reqModel)
	if err != nil {
		return wsMessageModel, err
	}

	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, contactModel.UUID)

	global.LOGGER.Info(fmt.Sprintf("消息请求路径:%s,消息请求体:%s", endpointUrl, string(bt)))
	respWrapper := httpClientUtils.PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, fmt.Errorf("SendMessage.HttpErr,StatusCode：%s,Body:%s", respWrapper.StatusCode, respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("SendMessage.WsErr,解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("SendMessage.WsErr,请求ws失败,记录：sendUrl: %s,请求体reqModel：%+v,响应体wsMessageModel: %+v", endpointUrl, reqModel, wsMessageModel)
		//global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	//messageId := ""
	//// 类型断言，确保 ResponseData 是 map[string]interface{}
	//if responseData, ok := wsMessageModel.ResponseData.(map[string]interface{}); ok {
	//	// 从 map 中获取 "id" 字段
	//	messageId = responseData["ID"].(string)
	//	global.LOGGER.Info(fmt.Sprintf("SendMessage,sendID:%s 消息ID:%s", reqSendMessageModel.UUID, messageId))
	//}

	global.LOGGER.Info(fmt.Sprintf("SendMessage,sendID:%s 消息发送成功: %s", reqSendMessageModel.UUID, wsMessageModel))
	return wsMessageModel, nil
}

func (svc *GreetService) addMessageToQueueV2(message model.MessageModel) {
	svc.mu.Lock()
	svc.messageQueue = append(svc.messageQueue, message)
	svc.mu.Unlock()
}

// 获取或刷新登录token
func (svc *GreetService) RefreshTokenV2() {
	waitSeconds := 10 * time.Second

__retry:
	global.LOGGER.Info("开始获取token...")
	var httpClientUtils = utils.HttpClientUtils{
		Header: map[string][]string{"token": {svc.GlobalTaskConfig.Token}},
	}

	reqJson := fmt.Sprintf(`{"username": "%s","password": "%s"}`, svc.GlobalTaskConfig.Username, svc.GlobalTaskConfig.Password)
	respWrapper := httpClientUtils.PostJson(UserLoginUrl, reqJson, HttpRequestTimeout)
	body := respWrapper.Body
	if respWrapper.StatusCode != http.StatusOK {
		global.LOGGER.Error(fmt.Sprintf("刷新token失败1:%v 重试...\n", body))
		time.Sleep(waitSeconds)
		goto __retry
	}

	if gjson.Get(body, "code").Int() != http.StatusOK {
		global.LOGGER.Error(fmt.Sprintf("刷新token失败2:%v 重试...\n", body))
		time.Sleep(waitSeconds)
		goto __retry
	}

	token := gjson.Get(body, "data.token").String()
	svc.GlobalTaskConfig.Token = token
	global.LOGGER.Info(fmt.Sprintf("获取到token:%v \n", token))

}

func (svc *GreetService) ProcessMessageQueueToTableV2() {
	ticker := time.NewTicker(global.SERVER_CONFIG.SystemConfig.BatchToTableTime * time.Second) //处理消息队列的时间
	defer ticker.Stop()
	for range ticker.C {
		svc.mu.Lock()
		if len(svc.messageQueue) == 0 {
			svc.mu.Unlock()
			continue // 跳过当前循环，避免重复解锁
		}

		batch := svc.messageQueue
		svc.messageQueue = make([]model.MessageModel, 0)
		svc.mu.Unlock()

		global.LOGGER.Info("批量保存招呼消息到表...")
		grouped := groupBySeatID(batch)
		//保存到消息表中 同一个seat_id的消息批量添加
		if err := messageService.AddBatchGrouped(grouped); err != nil {
			global.LOGGER.Error("messageService.AddBatchGrouped:" + err.Error())
			return
		}
	}
}

func (svc *GreetService) AddContacts(contactModel model.ContactModel, phoneNumbers []string) (resp.WsMessageModel, error) {

	var wsMessageModel resp.WsMessageModel
	// 1. 构建请求体
	requestModel := &req.WsAddContactsModel{}
	requestModel.CgiRequest.Contacts = phoneNumbers

	// 2. 转换为JSON
	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("转换JSON失败: %v", err))
		return wsMessageModel, err
	}
	sendUrl := WsAddContactsUrl
	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, contactModel.UUID)
	global.LOGGER.Info("消息请求体：" + string(bt))
	respWrapper := httpClientUtils.PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New("addContacts失败" + respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("addContacts解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("addContacts失败,记录：sendUrl: %s,请求体reqModel：%+v,响应体wsMessageModel: %+v", endpointUrl, requestModel, wsMessageModel)
		//global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	return wsMessageModel, nil
}

func (svc *GreetService) Initsession(uuid string, phone string) (resp.WsMessageModel, error) {

	var wsMessageModel resp.WsMessageModel
	// 1. 构建请求体
	requestModel := &req.Initsession{}
	requestModel.CgiRequest.ToUserId = phone

	// 2. 转换为JSON
	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("转换JSON失败: %v", err))
		return wsMessageModel, err
	}
	sendUrl := WsInitsession
	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	global.LOGGER.Info("Initsession消息请求体：" + string(bt))
	respWrapper := httpClientUtils.PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New("initsession失败" + respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("initsession解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("initsession失败,记录：sendUrl: %s,请求体reqModel：%s ,响应体wsMessageModel: %+v", endpointUrl, string(bt), wsMessageModel)
		//global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	return wsMessageModel, nil
}

func (svc *GreetService) HandleHijack(rds *redis.Client, uuid string, accType int, phone string) {
	rdsKey := "hijack:" + uuid
	err := redisUtils.Set(rds, rdsKey, accType, 0)
	if err != nil {
		global.LOGGER.Error("Greet.HandleHijack:" + err.Error())
		return
	}
	global.LOGGER.Info(fmt.Sprintf("hijack缓存设置成功,key: %s ,val :%s", rdsKey, accType))
	_, err = svc.Deletechat("QR"+uuid, phone)
	if err != nil {
		global.LOGGER.Error("Greet.HandleHijack:" + err.Error())
		return
	}
	global.LOGGER.Info(fmt.Sprintf("删除会话成功 uuid:%s ,accType:%s ", uuid, accType))
	if accType == 3 {
		_, err = svc.Archive("QR"+uuid, phone)
		if err != nil {
			global.LOGGER.Error("Greet.HandleHijack:" + err.Error())
			return
		}
		global.LOGGER.Info(fmt.Sprintf("归档成功 uuid:%s ,accType:%s ", uuid, accType))

		_, err = svc.MuteChat("QR"+uuid, phone)
		if err != nil {
			global.LOGGER.Error("Greet.HandleHijack:" + err.Error())
			return
		}
		global.LOGGER.Info(fmt.Sprintf("静音成功 uuid:%s ,accType:%s ", uuid, accType))
	}

}
func (svc *GreetService) Deletechat(uuid string, phone string) (resp.WsMessageModel, error) {

	var wsMessageModel resp.WsMessageModel
	// 1. 构建请求体
	requestModel := &req.Deletechat{}
	requestModel.CgiRequest.ToUserId = phone

	// 2. 转换为JSON
	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("转换JSON失败: %v", err))
		return wsMessageModel, err
	}
	sendUrl := WsDeteleChatUrl
	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	global.LOGGER.Info("Deletechat消息请求体：" + string(bt))
	respWrapper := httpClientUtils.PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New("deletechat失败" + respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("deletechat解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("deletechat失败,记录：sendUrl: %s,请求体reqModel：%s ,响应体wsMessageModel: %+v", endpointUrl, string(bt), wsMessageModel)
		//global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	return wsMessageModel, nil
}

func (svc *GreetService) Archive(uuid string, phone string) (resp.WsMessageModel, error) {

	var wsMessageModel resp.WsMessageModel
	// 1. 构建请求体
	requestModel := &req.Archive{}
	requestModel.CgiRequest.ToUserId = phone
	requestModel.CgiRequest.Achive = true

	// 2. 转换为JSON
	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("转换JSON失败: %v", err))
		return wsMessageModel, err
	}
	sendUrl := WsAchiveUrl
	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	global.LOGGER.Info("Archive消息请求体：" + string(bt))
	respWrapper := httpClientUtils.PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New("Achive失败" + respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("Achive解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("Achive失败,记录：sendUrl: %s,请求体reqModel：%s ,响应体wsMessageModel: %+v", endpointUrl, string(bt), wsMessageModel)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	return wsMessageModel, nil
}

func (svc *GreetService) MuteChat(uuid string, phone string) (resp.WsMessageModel, error) {

	var wsMessageModel resp.WsMessageModel
	// 1. 构建请求体
	requestModel := &req.MuteChat{}
	requestModel.CgiRequest.ToUserId = phone
	requestModel.CgiRequest.IsMute = true

	// 2. 转换为JSON
	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("转换JSON失败: %v", err))
		return wsMessageModel, err
	}
	sendUrl := WsMuteChatUrl
	//请求ws
	var httpClientUtils utils.HttpClientUtils
	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	global.LOGGER.Info("MuteChat消息请求体：" + string(bt))
	respWrapper := httpClientUtils.PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != http.StatusOK {
		return wsMessageModel, errors.New("MuteChat失败" + respWrapper.Body)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("MuteChat解析ws返回的结果失败：%s", err.Error())
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("MuteChat失败,记录：sendUrl: %s,请求体reqModel：%s ,响应体wsMessageModel: %+v", endpointUrl, string(bt), wsMessageModel)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	return wsMessageModel, nil
}
