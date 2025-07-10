package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/utils"
)

// ChatAccountClient 账号管理客户端
type ChatAccountClient struct {
}

// NewChatAccountClient 创建账号管理客户端
func NewChatAccountClient() *ChatAccountClient {
	return &ChatAccountClient{}
}

// ChatAccountResponse 账号管理服务响应（通用响应）
type ChatAccountResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// AssignAccountRequest 账号分配请求体
type AssignAccountRequest struct {
	SeatId int64  `json:"seatId"` // 座席ID
	Area   string `json:"area"`   // 区域
}

// AccountData 账号数据结构
type AccountData struct {
	AccountName int64 `json:"accountName"` // 账号名称
	AccountType int64 `json:"accountType"` // 账号类型
}

// AssignAccountResponse 账号分配响应
type AssignAccountResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *AccountData `json:"data"`
}

// AssignAccount 分配账号（使用Consul动态获取服务地址）
func (c *ChatAccountClient) AssignAccount(request *AssignAccountRequest) (*AccountData, error) {
	// 使用包级别的 consulUtils 实例
	chatAccountBaseURL, err := consulUtils.GetChatAccountServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("AssignAccount-获取ChatAccount服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatAccount服务地址: %v", err)
	}

	// 转换为JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("AssignAccount-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用全局配置中的URL路径）
	endpointURL := chatAccountBaseURL + global.ChatAccountAssignUrl
	reqInfo := fmt.Sprintf("AssignAccount-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(request))

	// 发送HTTP POST请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("AssignAccount-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response AssignAccountResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("AssignAccount-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 {
		errInfo := fmt.Sprintf("AssignAccount-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("AssignAccount-分配成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return response.Data, nil
}
