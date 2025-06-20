package model

type ContactModel struct {
	ID              string `gorm:"column:id;type:bigint;primary_key;auto_increment" json:"id"`                      // 对话id
	DisplayName     string `gorm:"column:display_name;type:varchar(255);not null" json:"displayName"`               // 昵称/备注 displayName (即修改备注的remark)
	Avatar          string `gorm:"column:avatar;type:varchar(500);default:null" json:"avatar"`                      // 头像地址，没有就留空
	Unread          int64  `gorm:"column:unread;type:bigint;not null;default:0" json:"unread"`                      // 未读消息的数量，收到新消息时 +1，调用聊天记录接口时清零
	PlatformAccID   int64  `gorm:"column:platform_acc_id;type:bigint;not null" json:"platformAccId"`                // 平台账号id
	PlatformAcc     string `gorm:"column:platform_acc;type:varchar(100);not null" json:"platformAcc"`               // 平台账号
	TargetAcc       string `gorm:"column:target_acc;type:varchar(100);not null" json:"targetAcc"`                   // 目标账号
	AllowReply      int8   `gorm:"column:allow_reply;type:tinyint;not null;default:0" json:"allowReply"`            // 打招呼是否回复，0或1
	IsOnline        int8   `gorm:"column:is_online;type:tinyint;not null;default:0" json:"isOnline"`                // 是否在线，0或1
	IsCustomerBlock int8   `gorm:"column:is_customer_block;type:tinyint;not null;default:0" json:"isCustomerBlock"` // 是否黑名单，0或1
	IsPinTop        int8   `gorm:"column:is_pin_top;type:tinyint;not null;default:0" json:"isPinTop"`               // 是否置顶，0或1
	LastSendTime    int64  `gorm:"column:last_send_time;type:bigint;not null" json:"lastSendTime"`                  // 最新消息时间(秒)，用于判断前端缓存聊天记录是否需要重新拉取
	GroupID         string `gorm:"column:group_id;type:varchar(255);not null;default:0" json:"groupID"`             // 分组id
	SeatID          int64  `gorm:"column:seat_id;type:bigint;not null;default:0" json:"seatID"`                     // 座席id
	CreateTime      int64  `gorm:"column:create_time;type:datetime;default:null" json:"createTime"`                 // 创建时间
	UpdateTime      int64  `gorm:"column:update_time;type:datetime;default:null" json:"updateTime"`                 // 更新时间
}

func (ContactModel) TableName() string {
	return "contact"
}
