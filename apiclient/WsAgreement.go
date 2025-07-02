package apiclient

import (
	"encoding/json"
	"fmt"

	"git.exclouds.org/ww/greet_ws/global/enum"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/model/req"
	"git.exclouds.org/ww/greet_ws/model/resp"
	"git.exclouds.org/ww/greet_ws/utils"
)

// SessionClient 只支持直连方式
// 归档会话相关API

const HttpRequestTimeout = 100

type WsAgreement struct{}

// DeleteChat 删除会话
func (c *WsAgreement) DeleteChat(uuid, phone string) (resp.WsMessageModel, error) {
	var wsMessageModel resp.WsMessageModel
	requestModel := &req.Deletechat{}
	requestModel.CgiRequest.ToUserId = phone

	bt, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("DeleteChat-转换JSON失败: %v", err))
		return wsMessageModel, err
	}

	endpointUrl := fmt.Sprintf("%s?uuid=%s", global.WsDeteleChatUrl, uuid)
	reqInfo := fmt.Sprintf("DeleteChat-请求URL: %s, 请求参数: %+v", endpointUrl, toJSONString(requestModel))

	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		errInfo := fmt.Sprintf("DeleteChat-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		errInfo := fmt.Sprintf("DeleteChat-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("DeleteChat-业务失败: %+v, 请求: %s", wsMessageModel, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("DeleteChat-: 删除会话: %s, 请求: %s", toJSONString(wsMessageModel), reqInfo))
	return wsMessageModel, nil
}

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

	endpointUrl := fmt.Sprintf("%s?uuid=%s", global.WsAchiveUrl, uuid)
	reqInfo := fmt.Sprintf("Archive-请求URL: %s, 请求参数: %+v", endpointUrl, toJSONString(requestModel))

	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		errInfo := fmt.Sprintf("Archive-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		errInfo := fmt.Sprintf("Archive-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("Archive-业务失败: %+v, 请求: %s", wsMessageModel, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("Archive-归档成功: 响应数据: %s, 请求: %s", toJSONString(wsMessageModel), reqInfo))
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

	endpointUrl := fmt.Sprintf("%s?uuid=%s", global.WsMuteChatUrl, uuid)
	reqInfo := fmt.Sprintf("MuteChat-请求URL: %s, 请求参数: %+v", endpointUrl, toJSONString(requestModel))

	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		errInfo := fmt.Sprintf("MuteChat-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		errInfo := fmt.Sprintf("MuteChat-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("MuteChat-业务失败: %+v, 请求: %s", wsMessageModel, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("MuteChat-静音成功: 响应数据: %s, 请求: %s", toJSONString(wsMessageModel), reqInfo))
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

	endpointUrl := fmt.Sprintf("%s?uuid=%s", global.WsInitsession, uuid)
	reqInfo := fmt.Sprintf("Initsession-请求URL: %s, 请求参数: %+v", endpointUrl, toJSONString(requestModel))

	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		errInfo := fmt.Sprintf("Initsession-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		errInfo := fmt.Sprintf("Initsession-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("Initsession-业务失败: %+v, 请求: %s", wsMessageModel, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("Initsession-初始化成功: 响应数据: %s, 请求: %s", toJSONString(wsMessageModel), reqInfo))
	return wsMessageModel, nil
}

// SendMessage 通用发送消息（支持文本、图片）
func (c *WsAgreement) SendMessage(reqModel interface{}, uuid string, messageType string) (resp.WsMessageModel, error) {
	var wsMessageModel resp.WsMessageModel
	sendUrl := ""
	switch messageType {
	case enum.MessageTypeText: //文字消息
		sendUrl = global.WsSendTxtUrl
	case enum.MessageTypeMedia: //图片消息
		sendUrl = global.WsSendImageUrl
	default:
		global.LOGGER.Error(fmt.Sprintf("SendMessage-消息类型不支持: %s", messageType))
		return wsMessageModel, fmt.Errorf("消息类型不支持：%s", messageType)
	}

	bt, err := json.Marshal(reqModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SendMessage-转换JSON失败: %v", err))
		return wsMessageModel, err
	}

	endpointUrl := fmt.Sprintf("%s?uuid=%s", sendUrl, uuid)
	reqInfo := fmt.Sprintf("SendMessage-请求URL: %s, 请求参数: %+v, 消息类型: %s", endpointUrl, toJSONString(reqModel), messageType)

	respWrapper := (&utils.HttpClientUtils{}).PostJson(endpointUrl, string(bt), HttpRequestTimeout)
	if respWrapper.StatusCode != 200 {
		errInfo := fmt.Sprintf("SendMessage-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if err = json.Unmarshal([]byte(respWrapper.Body), &wsMessageModel); err != nil {
		errInfo := fmt.Sprintf("SendMessage-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	if wsMessageModel.CgiBaseResponse.Ret != 0 {
		errInfo := fmt.Sprintf("SendMessage-业务失败: %+v, 请求: %s", wsMessageModel, reqInfo)
		global.LOGGER.Error(errInfo)
		return wsMessageModel, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("SendMessage-发送成功: 响应数据: %s, 请求: %s", toJSONString(wsMessageModel), reqInfo))
	return wsMessageModel, nil
}

// NewSessionClient 返回默认直连模式的SessionClient
func NewWsAgreement() *WsAgreement {
	return &WsAgreement{}
}
