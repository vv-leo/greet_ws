package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"git.exclouds.org/ww/greet_ws/global/enum"
	"git.exclouds.org/ww/greet_ws/model"
	"git.exclouds.org/ww/greet_ws/model/req"
	"git.exclouds.org/ww/greet_ws/model/resp"
	"git.exclouds.org/ww/greet_ws/utils"

	"git.exclouds.org/ww/greet_ws/apiclient"
	"git.exclouds.org/ww/greet_ws/global"
)

// 在包级别定义，整个应用共享一个实例
var (
	chatWorkClient    = apiclient.NewChatWorkClient()
	chatAccountClient = apiclient.NewChatAccountClient()
	chatHistClient    = apiclient.NewChatHistoryClient()
	wsAgreeClient     = apiclient.NewWsAgreement()
	rng               = rand.New(rand.NewSource(time.Now().UnixNano()))
	redisUtils        utils.RedisUtils
)

type GreetService struct {
	mu sync.Mutex
}

func NewGreetService() *GreetService {
	return &GreetService{}
}

func (svc *GreetService) GreetTask() {
	// 查询待处理的问候任务
	taskList, err := chatWorkClient.QueryWaitingTask()
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->查询待处理任务失败: %v", err))
		return
	}

	if len(taskList) == 0 {
		global.LOGGER.Info("greet->没有待处理的打招呼任务")
		return
	}

	global.LOGGER.Info(fmt.Sprintf("greet->查询到%d个租户待处理任务", len(taskList)))

	// 遍历处理每个租户的任务
	for _, taskData := range taskList {
		// 使用goroutine异步处理每个租户的任务
		go svc.handelTenantTask(taskData)
	}
}

func (svc *GreetService) handelTenantTask(taskData apiclient.TaskData) {
	tenantId := taskData.TenantId
	configRes, err := chatWorkClient.GetTenantHelloConfig(tenantId)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> 查询租户配置信息错误: %v", tenantId, err))
		return
	}
	if len(configRes) == 0 {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> 租户问候配置为空", tenantId))
		return
	}
	config := configRes[0]
	targets := taskData.Targets
	if len(targets) == 0 {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> 租户待处理任务为空", tenantId))
		return
	}

	for _, target := range targets {
		if target.TargetAcc == 0 || target.Country == "" || target.SeatId == "" {
			global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> 联系人数据有误跳过处理,target: %s", tenantId, toJSONString(target)))
			continue
		}
		// 使用goroutine异步处理，避免阻塞其他任务
		go svc.handTargetSendMess(target, config)
		time.Sleep(900 * time.Millisecond)
	}
}

func (svc *GreetService) handTargetSendMess(target apiclient.TaskTarget, config apiclient.TenantHelloConfigData) {
	tenantId := config.TenantId
	messageType := config.Type
	textList := config.ContentText
	content := ""
	messageId := ""
	// 将SeatId从string转换为int64
	seatId, err := strconv.ParseInt(target.SeatId, 10, 64)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> SeatId类型转换失败, err: %v", tenantId, target.SeatId, target.TargetAcc, err))
		return
	}

	// 创建账号分配请求
	pullAccReq := &apiclient.AssignAccountRequest{
		SeatId: seatId,
		Area:   target.Country,
	}

	platformAccInfo, err := chatAccountClient.AssignAccount(pullAccReq)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> 平台账号分配失败, err: %v", tenantId, target.SeatId, target.TargetAcc, err))
		return
	}

	platformAcc := platformAccInfo.AccountName
	accType := platformAccInfo.AccountType
	platformAccStr := strconv.FormatInt(platformAcc, 10)
	targetAccStr := strconv.FormatInt(target.TargetAcc, 10)
	uuid := platformAccStr
	if accType == 3 {
		uuid = "QR" + platformAccStr
	}

	if platformAcc == 0 || accType == 0 {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 平台账号数据错误, accountType: %d", tenantId, target.SeatId, target.TargetAcc, platformAcc, accType))
		return
	}

	contactGroup := fmt.Sprintf("%d:%d", platformAcc, target.TargetAcc)
	rds := global.REDIS
	contactGroupRdsKey := "contact_group:" + contactGroup
	exists, err := redisUtils.KeyExists(rds, contactGroupRdsKey)
	if exists {
		errMsg := fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 联系人组合已存在，跳过处理, contactGroup: %s", tenantId, target.SeatId, target.TargetAcc, platformAcc, contactGroup)
		global.LOGGER.Error(errMsg)

		status := "1"
		detailReq := &apiclient.ExecuteDetailRequest{
			TargetAcc:   &targetAccStr,
			SeatId:      &target.SeatId,
			PlatformAcc: &platformAccStr,
			Status:      &status,
			FailReson:   &errMsg,
		}
		_, err = chatWorkClient.ExecuteDetail(detailReq)
		if err != nil {
			global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 上报执行详情失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
			return
		}
		global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 上报执行详情成功", tenantId, target.SeatId, target.TargetAcc, platformAcc))
		return
	}

	//初始化会话
	_, err = wsAgreeClient.Initsession(uuid, strconv.FormatInt(platformAcc, 10))
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 初始化会话失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
		//TODO 上报
		return
	}

	cgiRequest := req.CgiRequest{
		ToUserId:  targetAccStr,
		IsGroup:   false,
		IsFakeMsg: false,
	}
	var wsMessModel = resp.WsMessageModel{}
	if messageType == "text" {
		content = getRandomGreet(textList)
		reqModel := req.WsSendMessageTextModel{
			WithCgiRequest: req.WithCgiRequest{
				CgiRequest:  cgiRequest,
				StanzaId:    "",
				Participant: "",
				Content:     content,
			},
		}
		//ws发消息，打招呼(文本)
		wsMessModel, err = wsAgreeClient.SendMessage(reqModel, uuid, enum.MessageTypeText)
		if err != nil {
			global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 文本消息发送失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
			return
		}
		global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 文本消息发送成功, message: %+v", tenantId, target.SeatId, target.TargetAcc, platformAcc, wsMessModel))
	} else if messageType == "bot" {
		//sourceUrlList := wsMsgModel.TaskContent.UrlList
		//reqModel := req.WsSendMessageRichTextModel{
		//	req.RichTextCgiRequest{
		//		ToUserId:  toUserId,
		//		IsGroup:   false,
		//		Title:     wsMsgModel.TaskContent.Title,
		//		Text:      wsMsgModel.TaskContent.SecondaryTitle,
		//		Body:      wsMsgModel.TaskContent.Content,
		//		ImageUrl:  wsMsgModel.TaskContent.ImageUrl,
		//		SourceUrl: getRandomGreet(sourceUrlList),
		//	},
		//}
		//wsSendMessageModel, err := svc.WsSendAndSaveRichTextMessageV2(reqModel, contactModel, req.SendMessageModel{Type: messageType}, tenantID, targetId, wsAccountId)
		//if err != nil {
		//	global.LOGGER.Error("图文===>打招呼失败：" + err.Error())
		//	_ = svc.finishAccount(tenantID, strconv.FormatInt(seatId, 10), wsAccountId, "", waAccount, "", "", wsSendMessageModel.CgiBaseResponse.ErrMsg)
		//	return err
		//}
		//global.LOGGER.Info("图文===>打招呼成功，结果：" + fmt.Sprintf("%+v", wsSendMessageModel))
	}
	if responseData, ok := wsMessModel.ResponseData.(map[string]interface{}); ok {
		// 从 map 中获取 "id" 字段
		messageId = responseData["ID"].(string)
		global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 消息ID: %s", tenantId, target.SeatId, target.TargetAcc, platformAcc, messageId))
	}
	contactID := global.GenerateUniqueId()
	// 创建联系人模型
	contactModel := model.ContactModel{
		Id:              contactID,
		DisplayName:     "",
		Avatar:          "",
		Unread:          0,
		PlatformAccID:   0,
		PlatformAcc:     platformAccStr,
		TargetAcc:       targetAccStr,
		AllowReply:      0,
		IsOnline:        1,
		IsCustomerBlock: 0,
		IsPinTop:        0,
		LastSendTime:    time.Now().Unix(),
		GroupId:         "",
		SeatId:          seatId, // 座席ID - 需要填充
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
	}

	// 保存联系人
	_, err = chatHistClient.SaveContact(contactModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 保存联系人失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
		return
	}

	// 创建消息模型
	messageModel := model.MessageModel{
		Id:          messageId,
		Status:      "0",
		ContactId:   contactID,
		SendTime:    time.Now().Unix(),
		DisplayName: "",
		Avatar:      "",
		IsFromMe:    1,
		Type:        messageType,
		Content:     content,
		SeatId:      seatId,
		MessageId:   messageId,
		TargetAcc:   targetAccStr,
		PlatformAcc: platformAccStr,
		CreateTime:  time.Now().Unix(),
		UpdateTime:  time.Now().Unix(),
	}

	// 保存消息
	_, err = chatHistClient.SaveMessage(messageModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 保存消息失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
		return
	}

	//处理劫持号
	if accType == 2 || accType == 3 {
		rdsKey := "hijack:" + platformAccStr
		err = redisUtils.Set(rds, rdsKey, accType, 0)
		if err != nil {
			global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> hijack缓存设置失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
			return
		}
		global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> hijack缓存设置成功, key: %s, val: %d", tenantId, target.SeatId, target.TargetAcc, platformAcc, rdsKey, accType))

		// 删除会话
		_, err = wsAgreeClient.DeleteChat(uuid, targetAccStr)
		if err != nil {
			global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 删除会话失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
			return
		}
		global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 删除会话成功, uuid: %s, accType: %d", tenantId, target.SeatId, target.TargetAcc, platformAcc, uuid, accType))

		if accType == 3 {
			// 归档会话
			_, err = wsAgreeClient.Archive(uuid, targetAccStr)
			if err != nil {
				global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 归档会话失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
				return
			}
			global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 归档会话成功, uuid: %s, accType: %d", tenantId, target.SeatId, target.TargetAcc, platformAcc, uuid, accType))

			// 静音会话
			_, err = wsAgreeClient.MuteChat(uuid, targetAccStr)
			if err != nil {
				global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 静音会话失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
				return
			}
			global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> 静音会话成功, uuid: %s, accType: %d", tenantId, target.SeatId, target.TargetAcc, platformAcc, uuid, accType))
		}
	}

	err = redisUtils.Set(rds, contactGroupRdsKey, seatId, time.Hour*24*8)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> contact_group缓存设置失败: %v", tenantId, target.SeatId, target.TargetAcc, platformAcc, err))
		return
	}
	global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> contact_group缓存设置成功, key:%s, val:%d", tenantId, target.SeatId, target.TargetAcc, platformAcc, contactGroupRdsKey, seatId))

	global.LOGGER.Info(fmt.Sprintf("greet->tenantId:%s -> seatId:%s -> targetAcc:%d -> platformAcc:%d -> greet任务执行成功", tenantId, target.SeatId, target.TargetAcc, platformAcc))
}

// // 随机取一条方法
func getRandomGreet(greetTextList []string) string {
	// 检查列表是否为空
	// 使用专用的随机数生成器
	randomIndex := rng.Intn(len(greetTextList))
	return greetTextList[randomIndex]
}

func toJSONString(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}
	return string(jsonData)
}
