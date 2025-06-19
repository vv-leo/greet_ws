package model

import (
	"time"
)

// TaskModel represents a task with detailed configuration.
type TaskModel struct {
	TaskID                 string             `json:"taskId"`                  // 任务唯一标识，用于跟踪任务状态
	MaxSessionCountPerUser int                `json:"maxSessionCountPerUser "` // 每个坐席允许创建的最大会话数量
	SubUserIdList          []string           `json:"subUserIdList"`           // 坐席id 列表
	Type                   string             `json:"type"`                    // 联系人数据类型 ('file' 或 'array')
	Contacts               interface{}        `json:"contacts"`                // 联系人手机号数组或联系人文件下载地址 (支持多类型数据)  //interface{}
	Messages               []TaskContentModel `json:"messages"`                // 发送的消息数组
	TaskConfig             TaskConfig         `json:"taskConfig"`              // 任务配置参数
	CreatedTime            time.Time          `json:"createdTime"`             // 任务创建时间，ISO 8601 格式 (如 '2024-11-17T08:00:00Z')
	Status                 string             `json:"status"`                  // 当前任务状态 ('pending', 'in-progress', 'completed', 'failed')
	Progress               float64            `json:"progress"`                // 任务进度，表示当前完成的百分比 (0 到 100)
}

// TaskContentModel represents a message to be sent.
type TaskContentModel struct {
	MessageType string `json:"messageType"` // 消息类型，例如：'text'
	Content     string `json:"content"`     // 消息内容
}

type PullAccountModel struct {
	Account             string `json:"account"`
	MessageCountPerRate string `json:"messageCountPerRate"` //这个ws账号允许加多少新联系人
}

type PullAccountModelV2 struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	NodeId   string `json:"nodeId"`
	Type     int    `json:"type"`
}

// TaskConfig represents configuration parameters for a task.
type TaskConfig struct {
	Priority               int                    `json:"priority,omitempty"`            // 任务优先级，可选值如 'high', 'medium', 'low'
	RetryCount             int                    `json:"retryCount,omitempty"`          // 最大重试次数，默认值根据业务需要设置
	SendIntervalSeconds    int64                  `json:"sendIntervalSeconds,omitempty"` // 每条消息之间的发送间隔时间 (单位：秒)
	AdditionalConfig       map[string]interface{} `json:"-"`                             // 扩展其他配置参数
	MessageCountPerAccount int                    `json:"messageCountPerAccount"`        // 表示一个ws账号分配多少新会话
}
