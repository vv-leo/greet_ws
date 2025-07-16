package model

type ContactModel struct {
	Id              string `json:"id"`              // 对话id
	DisplayName     string `json:"displayName"`     // 昵称/备注 displayName (即修改备注的remark)
	Avatar          string `json:"avatar"`          // 头像地址，没有就留空
	Unread          int64  `json:"unread"`          // 未读消息的数量，收到新消息时 +1，调用聊天记录接口时清零
	PlatformAccID   int64  `json:"platformAccId"`   // 平台账号id
	PlatformAcc     string `json:"platformAcc"`     // 平台账号
	PlatformAccType int8   `json:"platformAccType"` // 平台账号类型
	TargetAcc       string `json:"targetAcc"`       // 目标账号
	AllowReply      int8   `json:"allowReply"`      // 打招呼是否回复，0或1
	IsOnline        int8   `json:"isOnline"`        // 是否在线，0或1
	IsCustomerBlock int8   `json:"isCustomerBlock"` // 是否黑名单，0或1
	IsPinTop        int8   `json:"isPinTop"`        // 是否置顶，0或1
	LastSendTime    int64  `json:"lastSendTime"`    // 最新消息时间(秒)，用于判断前端缓存聊天记录是否需要重新拉取
	GroupId         string `json:"groupId"`         // 分组id
	SeatId          int64  `json:"seatId"`          // 座席id
	CreateTime      int64  `json:"createTime"`      // 创建时间
	UpdateTime      int64  `json:"updateTime"`      // 更新时间
}

func (ContactModel) TableName() string {
	return "contact"
}
