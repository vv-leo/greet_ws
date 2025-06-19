package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.exclouds.org/ww/greet_ws/global/enum"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/model/req"
	"git.exclouds.org/ww/greet_ws/model/resp"
	"git.exclouds.org/ww/greet_ws/service"
	"git.exclouds.org/ww/greet_ws/utils"
)

// SessionClient 只支持直连方式
// 归档会话相关API
// 如需DNS/Consul方式请新建独立文件

const HttpRequestTimeout = 100

type WsAgreement struct{}

// Archive 归档会话
func (c *WsAgreement) Archive(uuid, phone string) (resp.WsMessageModel, error) {
	var wsMessageModel resp.WsMessageModel
	requestModel := &req.Archive{}
	requestModel.CgiRequest.ToUserId = phone
	requestModel.CgiRequest.Achive = true

	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("Archive-转换JSON失败: %v", err))
		return wsMessageModel, err
	}

	endpointUrl := fmt.Sprintf("%s?uuid=%s", service.WsAchiveUrl, uuid)
	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		return wsMessageModel, errors.New("Archive失败" + respWrapper.Body)
	}
	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("Archive-解析ws返回的结果失败：%s", err.Error())
	}
	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		return wsMessageModel, fmt.Errorf("Archive失败,响应体: %+v", wsMessageModel)
	}
	return wsMessageModel, nil
}

// MuteChat 静音会话
func (c *WsAgreement) MuteChat(uuid, phone string) (resp.WsMessageModel, error) {
	var wsMessageModel resp.WsMessageModel
	requestModel := &req.MuteChat{}
	requestModel.CgiRequest.ToUserId = phone
	requestModel.CgiRequest.IsMute = true

	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("MuteChat-转换JSON失败: %v", err))
		return wsMessageModel, err
	}

	endpointUrl := fmt.Sprintf("%s?uuid=%s", service.WsMuteChatUrl, uuid)
	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		return wsMessageModel, errors.New("MuteChat失败" + respWrapper.Body)
	}
	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("MuteChat-解析ws返回的结果失败：%s", err.Error())
	}
	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		return wsMessageModel, fmt.Errorf("MuteChat失败,响应体: %+v", wsMessageModel)
	}
	return wsMessageModel, nil
}

// Initsession 初始化会话
func (c *WsAgreement) Initsession(uuid, phone string) (resp.WsMessageModel, error) {
	var wsMessageModel resp.WsMessageModel
	requestModel := &req.Initsession{}
	requestModel.CgiRequest.ToUserId = phone

	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("Initsession-转换JSON失败: %v", err))
		return wsMessageModel, err
	}

	endpointUrl := fmt.Sprintf("%s?uuid=%s", service.WsInitsession, uuid)
	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		return wsMessageModel, errors.New("Initsession失败" + respWrapper.Body)
	}
	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("Initsession-解析ws返回的结果失败：%s", err.Error())
	}
	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		return wsMessageModel, fmt.Errorf("Initsession失败,响应体: %+v", wsMessageModel)
	}
	return wsMessageModel, nil
}

// SendMessage 通用发送消息（支持文本、图片）
func (c *WsAgreement) SendMessage(reqModel interface{}, uuid string, messageType string) (resp.WsMessageModel, error) {
	var wsMessageModel resp.WsMessageModel
	sendUrl := ""
	switch messageType {
	case enum.MessageTypeText: //文字消息
		sendUrl = service.WsSendTxtUrl
	case enum.MessageTypeMedia: //图片消息
		sendUrl = service.WsSendImageUrl
	default:
		return wsMessageModel, errors.New("消息类型不支持：" + messageType)
	}

	bt, err := json.Marshal(reqModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SendMessage-转换JSON失败: %v", err))
		return wsMessageModel, err
	}

	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		return wsMessageModel, errors.New("SendMessage失败" + respWrapper.Body)
	}
	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		return wsMessageModel, fmt.Errorf("SendMessage-解析ws返回的结果失败：%s", err.Error())
	}
	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		return wsMessageModel, fmt.Errorf("SendMessage失败,响应体: %+v", wsMessageModel)
	}
	global.LOGGER.Info(fmt.Sprintf("SendMessage,uuid:%s, type:%s, 消息发送成功: %+v", uuid, messageType, wsMessageModel))
	return wsMessageModel, nil
}

// NewSessionClient 返回默认直连模式的SessionClient
func NewWsAgreement() *WsAgreement {
	return &WsAgreement{}
}
