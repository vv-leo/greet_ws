package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/utils"
)

// ChatWorkClient 工作任务客户端
type ChatWorkClient struct {
}

// NewChatWorkClient 创建工作任务客户端
func NewChatWorkClient() *ChatWorkClient {
	return &ChatWorkClient{}
}

// ChatWorkResponse 工作任务服务响应（通用响应）
type ChatWorkResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TaskTarget 任务目标结构
type TaskTarget struct {
	TargetAcc int64  `json:"targetAcc"`
	Country   string `json:"country"`
	SeatId    string `json:"seatId"`
}

// TaskData 任务数据结构
type TaskData struct {
	TenantId string       `json:"tenantId"`
	Targets  []TaskTarget `json:"targets"`
}

// QueryWaitingTaskResponse 查询等待任务响应
type QueryWaitingTaskResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []TaskData `json:"data"`
}

// ExecuteDetailRequest 执行任务详情请求体
type ExecuteDetailRequest struct {
	TargetAcc   *string `json:"targetAcc"`   // 目标账号
	PlatformAcc *string `json:"platformAcc"` // 平台账号
	SeatId      *string `json:"seatId"`      // 座席ID
	Status      *string `json:"status"`      // 状态
	FailReson   *string `json:"failReson"`   // 失败原因
}

// TenantHelloConfigData 租户问候配置数据结构
type TenantHelloConfigData struct {
	TenantId    string      `json:"tenantId"`    // 租户ID
	Type        string      `json:"type"`        // 类型：text|link
	Content     []string    `json:"content"`     // 文本内容
	ContentLink interface{} `json:"contentLink"` // 链接内容
}

// GetTenantHelloConfigRequest 获取租户问候配置请求体
type GetTenantHelloConfigRequest struct {
	TenantId string `json:"tenantId"` // 租户ID
}

// GetTenantHelloConfigResponse 获取租户问候配置响应
type GetTenantHelloConfigResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    []TenantHelloConfigData `json:"data"`
}

// QueryWaitingTask 查询等待任务（使用Consul动态获取服务地址）
func (c *ChatWorkClient) QueryWaitingTask() ([]TaskData, error) {
	// 使用包级别的 consulUtils 实例
	chatWorkBaseURL, err := consulUtils.GetChatWorkServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("QueryWaitingTask-获取ChatWork服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatWork服务地址: %v", err)
	}

	// 构建完整的endpoint URL（使用全局配置中的URL路径）
	endpointURL := chatWorkBaseURL + global.ChatWorkQueryWaitingTaskUrl
	reqInfo := fmt.Sprintf("QueryWaitingTask-请求URL: %s", endpointURL)

	// 发送HTTP POST请求（入参为空JSON对象）
	respWrapper := httpClient.PostJson(endpointURL, "{}", DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("QueryWaitingTask-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response QueryWaitingTaskResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("QueryWaitingTask-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 {
		errInfo := fmt.Sprintf("QueryWaitingTask-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("QueryWaitingTask-查询成功: 响应数据: %s, 请求: %s", toJSONString(response.Data), reqInfo))
	return response.Data, nil
}

// ExecuteDetail 执行任务详情（使用Consul动态获取服务地址）
func (c *ChatWorkClient) ExecuteDetail(request *ExecuteDetailRequest) (*ChatWorkResponse, error) {
	// 使用包级别的 consulUtils 实例
	chatWorkBaseURL, err := consulUtils.GetChatWorkServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("ExecuteDetail-获取ChatWork服务地址失败: %v", err))
		return nil, fmt.Errorf("无法获取ChatWork服务地址: %v", err)
	}

	// 转换为JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("ExecuteDetail-转换JSON失败: %v", err))
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用全局配置中的URL路径）
	endpointURL := chatWorkBaseURL + global.ChatWorkExecuteDetailUrl
	reqInfo := fmt.Sprintf("ExecuteDetail-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(request))

	// 发送HTTP POST请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("ExecuteDetail-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response ChatWorkResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("ExecuteDetail-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return nil, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 {
		errInfo := fmt.Sprintf("ExecuteDetail-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return &response, fmt.Errorf(errInfo)
	}

	global.LOGGER.Info(fmt.Sprintf("ExecuteDetail-执行成功: 响应数据: %s, 请求: %s", toJSONString(response), reqInfo))
	return &response, nil
}

// GetTenantHelloConfig 获取租户问候配置（使用Consul动态获取服务地址）
func (c *ChatWorkClient) GetTenantHelloConfig(tenantId string) (TenantHelloConfigData, error) {
	var emptyConfig TenantHelloConfigData

	// 使用包级别的 consulUtils 实例
	chatWorkBaseURL, err := consulUtils.GetChatWorkServiceURLWithStrategy(utils.RoundRobin)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("GetTenantHelloConfig-获取ChatWork服务地址失败: %v", err))
		return emptyConfig, fmt.Errorf("无法获取ChatWork服务地址: %v", err)
	}

	// 构建请求体
	request := GetTenantHelloConfigRequest{
		TenantId: tenantId,
	}

	// 转换为JSON
	requestData, err := json.Marshal(request)
	if err != nil {
		global.LOGGER.Error(fmt.Sprintf("GetTenantHelloConfig-转换JSON失败: %v", err))
		return emptyConfig, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	// 构建完整的endpoint URL（使用全局配置中的URL路径）
	endpointURL := chatWorkBaseURL + global.ChatWorkTenantHelloUrl
	reqInfo := fmt.Sprintf("GetTenantHelloConfig-请求URL: %s, 请求参数: %+v", endpointURL, toJSONString(request))

	// 发送HTTP POST请求
	respWrapper := httpClient.PostJson(endpointURL, string(requestData), DefaultRequestTimeout)

	// 检查HTTP状态码
	if respWrapper.StatusCode != http.StatusOK {
		errInfo := fmt.Sprintf("GetTenantHelloConfig-HTTP请求失败: 状态码%d, 响应: %s, 请求: %s", respWrapper.StatusCode, respWrapper.Body, reqInfo)
		global.LOGGER.Error(errInfo)
		return emptyConfig, fmt.Errorf(errInfo)
	}

	// 解析响应
	var response GetTenantHelloConfigResponse
	if err := json.Unmarshal([]byte(respWrapper.Body), &response); err != nil {
		errInfo := fmt.Sprintf("GetTenantHelloConfig-解析响应失败: %v, 请求: %s", err, reqInfo)
		global.LOGGER.Error(errInfo)
		return emptyConfig, fmt.Errorf(errInfo)
	}

	// 检查业务响应码
	if response.Code != 200 {
		errInfo := fmt.Sprintf("GetTenantHelloConfig-业务失败: %+v, 请求: %s", response, reqInfo)
		global.LOGGER.Error(errInfo)
		return emptyConfig, fmt.Errorf(errInfo)
	}

	// 从数组中取第一个元素返回
	if len(response.Data) > 0 {
		global.LOGGER.Info(fmt.Sprintf("GetTenantHelloConfig-获取成功: 响应数据: %s, 请求: %s", toJSONString(response.Data[0]), reqInfo))
		return response.Data[0], nil
	}

	global.LOGGER.Info(fmt.Sprintf("GetTenantHelloConfig-获取成功但数据为空: 请求: %s", reqInfo))
	return emptyConfig, nil
}

// ReportExecuteDetail 上报执行详情的封装方法
func (c *ChatWorkClient) ReportExecuteDetail(targetAcc int64, seatId string, platformAcc int64, status string, failReason string) error {
	// 转换数字为字符串
	targetAccStr := strconv.FormatInt(targetAcc, 10)
	platformAccStr := strconv.FormatInt(platformAcc, 10)

	detailReq := &ExecuteDetailRequest{
		TargetAcc:   &targetAccStr,
		SeatId:      &seatId,
		PlatformAcc: &platformAccStr,
		Status:      &status,
	}
	if failReason != "" {
		detailReq.FailReson = &failReason
	}

	_, err := c.ExecuteDetail(detailReq)
	return err
}
