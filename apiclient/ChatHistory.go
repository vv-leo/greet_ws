package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/model"
	"git.exclouds.org/ww/greet_ws/utils"
)

// 常量定义
const (
	DefaultRequestTimeout = 100 // 超时时间（秒）
)

// 包级别的工具实例，避免重复创建
var (
	consulUtils = &utils.ConsulUtils{}
	httpClient  = &utils.HttpClientUtils{}
)

// toJSONString 将结构体转换为可读的JSON字符串
func toJSONString(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}
	return string(jsonData)
}

// ChatHistoryClient 聊天历史客户端
type ChatHistoryClient struct {
	baseURL string
}

// ChatHistoryResponse 聊天历史服务响应
type ChatHistoryResponse struct {
	Code    int         `json:"code"`
	Status  int         `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// GetContactRequest 获取联系人请求体
type GetContactRequest struct {
	Id          *string `json:"id"`
	SeatId      *int64  `json:"seatId"`
	PlatformAcc *string `json:"platformAcc"`
	TargetAcc   *string `json:"targetAcc"`
}

// GetContactResponse 获取联系人响应体
type GetContactResponse struct {
	Code    int                 `json:"code"`
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    *model.ContactModel `json:"data"`
}

// UpdateContactRequest 更新联系人请求体
type UpdateContactModel struct {
	ID           string `json:"id"`           // 对话id
	LastSendTime int64  `json:"lastSendTime"` // 最新消息时间(秒)，用于判断前端缓存聊天记录是否需要重新拉取
	SeatId       int64  `json:"seatID"`       // 座席id
	Unread       int64  `json:"unread"`       //未读
}

// GetMessagesByIdsModel 批量获取消息请求模型
type GetMessagesByIdsModel struct {
	MessageIds []string `json:"messageIds"`
	SeatId     int64    `json:"seatId"`
}

// UpdateMessageStatusModel 批量更新消息状态请求模型
type UpdateMessageStatusModel struct {
	MessageIds []string `json:"messageIds"`
	Status     string   `json:"status"`
	SeatId     int64    `json:"seatId"`
}

// NewChatHistoryClient 创建聊天历史客户端
func NewChatHistoryClient() *ChatHistoryClient {
	return &ChatHistoryClient{}
}

// SaveMessages 保存消息到聊天历史服务（使用Consul动态获取服务地址）
func (c *ChatHistoryClient) SaveMessage(message model.MessageModel) (*ChatHistoryResponse, error) {
	// 使用包级别的 consulUtils 实例
	chatHistoryBaseURL, err := consulUtils.GetChatHistoryServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SaveMessage-获取ChatHistory服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatHistory服务地址: %v", err)
	}
	// 转换为JSON
	requestData, err := json.Marshal(message)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SaveMessage-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用BaseService中的URL配置）
	endpointURL := chatHistoryBaseURL + global.ChatHistorySaveMessagesUrl
	reqInfo := fmt.Sprintf("SaveMessage-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(message))

	// 发送HTTP请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("SaveMessage-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response ChatHistoryResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("SaveMessage-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 || response.Status != 1 {
		errInfo := fmt.Sprintf("SaveMessage-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("SaveMessage-保存成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}

// SaveContacts 保存联系人到聊天历史服务（使用Consul动态获取服务地址）
func (c *ChatHistoryClient) SaveContact(contact model.ContactModel) (*ChatHistoryResponse, error) {

	// 使用包级别的 consulUtils 实例
	chatHistoryBaseURL, err := consulUtils.GetChatHistoryServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SaveContact-获取ChatHistory服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatHistory服务地址: %v", err)
	}
	// 转换为JSON
	requestData, err := json.Marshal(contact)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SaveContact-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用BaseService中的URL配置）
	endpointURL := chatHistoryBaseURL + global.ChatHistorySaveContactsUrl
	reqInfo := fmt.Sprintf("SaveContact-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(contact))

	// 发送HTTP请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("SaveContact-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response ChatHistoryResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("SaveContact-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 || response.Status != 1 {
		errInfo := fmt.Sprintf("SaveContact-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("SaveContact-保存成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}

// GetContact 从聊天历史服务获取联系人信息（使用Consul动态获取服务地址）
func (c *ChatHistoryClient) GetContact(contact *GetContactRequest) (*GetContactResponse, error) {
	if contact.SeatId == nil {
		global.LOGGER.Error("GetContact-坐席id不能为空")
		return nil, errors.New("坐席id不能为空")
	}

	// 使用包级别的 consulUtils 实例
	chatHistoryBaseURL, err := consulUtils.GetChatHistoryServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("GetContact-获取ChatHistory服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatHistory服务地址: %v", err)
	}
	// 转换为JSON
	requestData, err := json.Marshal(contact)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("GetContact-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用BaseService中的URL配置）
	endpointURL := chatHistoryBaseURL + global.ChatHistoryGetContactUrl
	reqInfo := fmt.Sprintf("GetContact-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(contact))

	// 发送HTTP请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("GetContact-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response GetContactResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("GetContact-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 || response.Status != 1 {
		errInfo := fmt.Sprintf("GetContact-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	if response.Data == nil {
		errInfo := fmt.Sprintf("GetContact-联系人不存在: %s, 请求: %s", response.Message, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("GetContact-获取成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}

// UpdateContact 更新联系人信息（使用Consul动态获取服务地址）
func (c *ChatHistoryClient) UpdateContact(updateContact UpdateContactModel) (*ChatHistoryResponse, error) {
	if updateContact.ID == "" {
		global.LOGGER.Error("UpdateContact-联系人ID不能为空")
		return nil, errors.New("联系人ID不能为空")
	}

	// 使用包级别的 consulUtils 实例
	chatHistoryBaseURL, err := consulUtils.GetChatHistoryServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("UpdateContact-获取ChatHistory服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatHistory服务地址: %v", err)
	}

	// 转换为JSON
	requestData, err := json.Marshal(updateContact)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("UpdateContact-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用BaseService中的URL配置）
	endpointURL := chatHistoryBaseURL + global.ChatHistoryUpdateContactUrl
	reqInfo := fmt.Sprintf("UpdateContact-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(updateContact))

	// 发送HTTP请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("UpdateContact-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response ChatHistoryResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("UpdateContact-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 || response.Status != 1 {
		errInfo := fmt.Sprintf("UpdateContact-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("UpdateContact-更新成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}

// GetMessagesByIds 批量获取消息（使用Consul动态获取服务地址）
func (c *ChatHistoryClient) GetMessagesByIds(request GetMessagesByIdsModel) (*ChatHistoryResponse, error) {
	// 使用包级别的 consulUtils 实例
	chatHistoryBaseURL, err := consulUtils.GetChatHistoryServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("GetMessagesByIds-获取ChatHistory服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatHistory服务地址: %v", err)
	}

	// 转换为JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("GetMessagesByIds-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL
	endpointURL := chatHistoryBaseURL + "/api/v1/im/getMessagesByIds"
	reqInfo := fmt.Sprintf("GetMessagesByIds-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(request))

	// 发送HTTP请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("GetMessagesByIds-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response ChatHistoryResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("GetMessagesByIds-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 || response.Status != 1 {
		errInfo := fmt.Sprintf("GetMessagesByIds-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("GetMessagesByIds-获取成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}

// UpdateMessageStatus 批量更新消息状态（使用Consul动态获取服务地址）
func (c *ChatHistoryClient) UpdateMessageStatus(request UpdateMessageStatusModel) (*ChatHistoryResponse, error) {
	// 使用包级别的 consulUtils 实例
	chatHistoryBaseURL, err := consulUtils.GetChatHistoryServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("UpdateMessageStatus-获取ChatHistory服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatHistory服务地址: %v", err)
	}

	// 转换为JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("UpdateMessageStatus-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL
	endpointURL := chatHistoryBaseURL + "/api/v1/im/updateMessageStatus"
	reqInfo := fmt.Sprintf("UpdateMessageStatus-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(request))

	// 发送HTTP请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("UpdateMessageStatus-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response ChatHistoryResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("UpdateMessageStatus-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 || response.Status != 1 {
		errInfo := fmt.Sprintf("UpdateMessageStatus-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("UpdateMessageStatus-更新成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}
