package model

type MessageModel struct {
	ID                string `gorm:"column:id;type:varchar(255);primary_key" json:"id"`                 // 消息的唯一id
	Status            string `gorm:"column:status;type:varchar(255);default:null" json:"status"`        // 消息状态
	ContactID         string `gorm:"column:contact_id;type:varchar(255);not null" json:"contactId"`     // 同对话id
	SendTime          int64  `gorm:"column:send_time;type:bigint;default:null" json:"sendTime"`         // 发送时间的时间戳
	TargetDisplayName string `gorm:"column:target_display_name;type:varchar(255);not null" json:"-"`    // 昵称/备注
	TargetAvatar      string `gorm:"column:target_avatar;type:varchar(255);not null" json:"-"`          // 头像地址 没有就留空就行
	IsFromMe          int8   `gorm:"column:is_from_me;type:tinyint;default:0" json:"is_from_me"`        // 自己发送的为1，联系人发送的为0
	Type              string `gorm:"column:type;type:varchar(50);default:null" json:"type"`             // 消息类型（如"event"、"image"、"text"）
	Content           string `gorm:"column:content;type:text" json:"content"`                           // 消息内容
	SeatID            int64  `gorm:"column:seat_id;type:bigint;not null;default:0" json:"-"`            // 座席id
	MessageID         string `gorm:"column:message_id;type:varchar(50);default:null" json:"-"`          // 来自的messageID
	TargetAcc         string `gorm:"column:target_acc;type:varchar(50);not null;default:''" json:"-"`   // 目标账号
	PlatformAcc       string `gorm:"column:platform_acc;type:varchar(50);not null;default:''" json:"-"` // 平台账号
	CreateTime        int64  `gorm:"column:create_time;type:datetime;default:null" json:"-"`            // 创建时间
	UpdateTime        int64  `gorm:"column:update_time;type:datetime;default:null" json:"-"`            // 更新时间
}

type ContentRichTextModel struct {
	Text        string `json:"text"`
	Type        string `json:"type"`
	ContextInfo struct {
	} `json:"contextInfo"`
}

// tableName 设置表名
func (MessageModel) TableName() string {
	return "message"
}
