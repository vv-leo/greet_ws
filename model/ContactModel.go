package model

type ContactModel struct {
	ID              string `gorm:"column:id;type:bigint;primary_key;auto_increment" json:"id"`                        // 对话id
	DisplayName     string `gorm:"column:display_name;type:varchar(255);not null" json:"displayName"`                 // 昵称/备注 displayName (即修改备注的remark)
	Avatar          string `gorm:"column:avatar;type:varchar(500);default:null" json:"avatar"`                        // 头像地址，没有就留空
	Index           string `gorm:"column:index;type:varchar(255);not null;default:'0'" json:"index"`                  // 这个就是index，留0
	Unread          int64  `gorm:"column:unread;type:bigint;not null;default:0" json:"unread"`                        // 未读消息的数量，收到新消息时 +1，调用聊天记录接口时清零
	Phone           string `gorm:"column:phone;type:varchar(15);not null" json:"phone"`                               // 对话的联系人手机号
	ContactID       string `gorm:"column:contact_id;type:bigint;not null;unique" json:"contact_id"`                   // 对话id
	WaAccountID     int64  `gorm:"column:wa_account_id;type:bigint;not null" json:"wa_account_id"`                    // ws账号id，系统生成的wa账号操作id
	WaAccount       string `gorm:"column:wa_account;type:varchar(100);not null" json:"-"`                             // 自己的ws账号
	WaAccountTo     string `gorm:"column:wa_account_to;type:varchar(100);not null" json:"-"`                          // 对方的ws账号
	AllowReply      int8   `gorm:"column:allow_reply;type:tinyint;not null;default:0" json:"allow_reply"`             // 打招呼是否回复，0或1
	IsOnline        int8   `gorm:"column:is_online;type:tinyint;not null;default:0" json:"is_online"`                 // 是否在线，0或1
	IsCustomerBlock int8   `gorm:"column:is_customer_block;type:tinyint;not null;default:0" json:"is_customer_block"` // 是否黑名单，0或1
	IsPinTop        int8   `gorm:"column:is_pin_top;type:tinyint;not null;default:0" json:"is_pin_top"`               // 是否置顶，0或1
	//ClientVer       int8   `gorm:"column:client_ver;type:tinyint;not null;default:1" json:"client_ver"`                 // 接口版本，默认填1
	//CanUseCallOffer int8   `gorm:"column:can_use_call_offer;type:tinyint;not null;default:0" json:"can_use_call_offer"` //
	LastSendTime int64  `gorm:"column:last_send_time;type:bigint;not null" json:"lastSendTime"` // 最新消息时间(秒)，用于判断前端缓存聊天记录是否需要重新拉取
	GroupID      string `gorm:"column:group_id;type:varchar(255);not null;default:0" json:"-"`  // 分组id
	UUID         string `gorm:"column:uuid;type:varchar(100);not null" json:"-"`                // 对列中来，用于ws发消息
	ToUserId     string `gorm:"column:to_user_id;type:varchar(100);not null" json:"-"`          // 对列中来ToUserId，用于ws发消息
	SeatID       int64  `gorm:"column:seat_id;type:bigint;not null;default:0" json:"-"`         // 座席id
	CreateTime   int64  `gorm:"column:create_time;type:datetime;default:null" json:"-"`         // 创建时间
	UpdateTime   int64  `gorm:"column:update_time;type:datetime;default:null" json:"-"`         // 更新时间
	NodeID       string `gorm:"column:node_id;type:varchar(255);default:not null" json:"-"`     // nodeId
}

func (ContactModel) TableName() string {
	return "contact"
}
