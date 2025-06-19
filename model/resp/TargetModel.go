package resp

type TargetModel struct {
	ID                string `json:"id"`                // 用户 ID
	Phone             string `json:"phone"`             // 电话号码
	Name              string `json:"name"`              // 姓名
	CountryCode       string `json:"countryCode"`       // 国家代码
	BelongToSubUserID string `json:"belongToSubUserId"` // 所属子用户 ID
	Assigned          bool   `json:"assigned"`          // 是否分配
	Sent              bool   `json:"sent"`              // 是否发送
	Doing             bool   `json:"doing"`             // 是否进行中
	AssignedTimestamp *int64 `json:"assignedTimestamp"` // 分配时间戳，允许为 null
	SentTimestamp     *int64 `json:"sentTimestamp"`     // 发送时间戳，允许为 null
	CreateTime        string `json:"createTime"`        // 创建时间
}
