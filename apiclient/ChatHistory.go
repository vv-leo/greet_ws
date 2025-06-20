package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/utils"
)

// 常量定义
const (
	DefaultRequestTimeout = 30 // 30秒
)

// ChatHistoryClient 聊天历史客户端
type ChatHistoryClient struct {
	baseURL string
}

// ChatHistoryMessageModel 聊天历史消息模型
type ChatHistoryMessageModel struct {
	ID         string `json:"id"`         // 消息唯一ID
	ContactID  string `json:"contactId"`  // 联系人ID
	Content    string `json:"content"`    // 消息内容
	Type       string `json:"type"`       // 消息类型（text、image、audio等）
	IsFromMe   int8   `json:"isFromMe"`   // 是否是自己发送的消息 1-是 0-否
	SendTime   int64  `json:"sendTime"`   // 发送时间戳
	Status     string `json:"status"`     // 消息状态
	SeatID     int64  `json:"seatId"`     // 座席ID
	MessageID  string `json:"messageId"`  // 消息ID
	FromJID    string `json:"fromJid"`    // 发送方JID
	UserJID    string `json:"userJid"`    // 用户JID
	CreateTime int64  `json:"createTime"` // 创建时间
	UpdateTime int64  `json:"updateTime"` // 更新时间
}

// SaveMessagesRequest 保存消息请求体
type SaveMessagesRequest struct {
	Messages []ChatHistoryMessageModel `json:"messages"`
}

// ChatHistoryResponse 聊天历史服务响应
type ChatHistoryResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// QueryMessagesRequest 查询消息请求体
type QueryMessagesRequest struct {
	ContactID   string `json:"contactId"`             // 联系人ID
	SeatID      int64  `json:"seatId,omitempty"`      // 座席ID
	Page        int    `json:"page"`                  // 页码
	Limit       int    `json:"limit"`                 // 每页数量
	StartTime   int64  `json:"startTime,omitempty"`   // 开始时间
	EndTime     int64  `json:"endTime,omitempty"`     // 结束时间
	MessageType string `json:"messageType,omitempty"` // 消息类型
}

// QueryMessagesResponse 查询消息响应
type QueryMessagesResponse struct {
	ChatHistoryResponse
	Data struct {
		Total    int64                     `json:"total"`
		Messages []ChatHistoryMessageModel `json:"messages"`
	} `json:"data"`
}

// DeleteMessagesRequest 删除消息请求体
type DeleteMessagesRequest struct {
	ContactID  string   `json:"contactId"`            // 联系人ID
	MessageIDs []string `json:"messageIds,omitempty"` // 消息ID列表，为空则删除该联系人所有消息
	SeatID     int64    `json:"seatId,omitempty"`     // 座席ID
}

// NewChatHistoryClient 创建聊天历史客户端
func NewChatHistoryClient(baseURL string) *ChatHistoryClient {
	return &ChatHistoryClient{
		baseURL: baseURL,
	}
}

// SaveMessages 保存消息到聊天历史服务
func (c *ChatHistoryClient) SaveMessages(messages []ChatHistoryMessageModel) (*ChatHistoryResponse, error) {
	if len(messages) == 0 {
		return nil, errors.New("消息列表不能为空")
	}

	// 构建请求体
	requestModel := &SaveMessagesRequest{
		Messages: messages,
	}

	// 转换为JSON
	requestData, err := json.Marshal(requestModel)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("SaveMessages-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL
	endpointURL := fmt.Sprintf("%s/api/v1/im/saveMessages", c.baseURL)

	global.LOGGER.Info(fmt.Sprintf("SaveMessages-请求URL: %s", endpointURL))
	global.LOGGER.Info(fmt.Sprintf("SaveMessages-请求体: %s", string(requestData)))

	// 发送HTTP请求
	httpClient := &utils.HttpClientUtils{}
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		global.LOGGER.Error(fmt.Sprintf("SaveMessages-HTTP请求失败: 状态码%d, 响应: %s",
			respWrapper.StatusCode, respWrapper.Body))
		return nil, fmt.Errorf("HTTP请求失败: 状态码%d, 响应: %s",
			respWrapper.StatusCode, respWrapper.Body)
	}

	// 解析响应
	var response ChatHistoryResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		global.LOGGER.Error(fmt.Sprintf("SaveMessages-解析响应失败: %v, 响应体: %s", err, respWrapper.Body))
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查业务响应码
	if response.Code != 200 && response.Code != 0 {
		global.LOGGER.Error(fmt.Sprintf("SaveMessages-业务失败: %+v", response))
		return &response, fmt.Errorf("保存消息失败: %s", response.Message)
	}

	global.LOGGER.Info(fmt.Sprintf("SaveMessages-保存成功: 共%d条消息", len(messages)))
	return &response, nil
}
