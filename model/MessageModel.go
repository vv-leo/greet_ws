package model

type MessageModel struct {
	Ack                  int8   `gorm:"column:ack;type:tinyint;not null;default:0" json:"ack"`              // 标识是否已读的标志，暂时固定为1
	WaAccountID          int64  `gorm:"column:wa_account_id;type:bigint;not null" json:"wa_account_id"`     // 发送消息的WS账号在系统中的id编号
	To                   string `gorm:"column:to;type:varchar(255);not null" json:"to"`                     // 对话id
	ID                   string `gorm:"column:id;type:varchar(255);primary_key" json:"id"`                  // 消息的唯一id
	Status               string `gorm:"column:status;type:varchar(255);default:null" json:"status"`         // 消息状态
	ToContactID          string `gorm:"column:to_contact_id;type:varchar(255);not null" json:"toContactId"` // 同对话id
	SendTime             int64  `gorm:"column:send_time;type:bigint;default:null" json:"sendTime"`          // 发送时间的时间戳
	FromUser_ID          string `gorm:"column:from_user_id;type:varchar(255);not null" json:"-"`            // 如果是自己发送的这里是1 如果是联系人发送的这里显示对话id
	FromUser_DisplayName string `gorm:"column:from_user_display_name;type:varchar(255);not null" json:"-"`  // 昵称/备注
	FromUser_Avatar      string `gorm:"column:from_user_avatar;type:varchar(255);not null" json:"-"`        // 头像地址 没有就留空就行
	IsFromMe             int8   `gorm:"column:is_from_me;type:tinyint;default:0" json:"is_from_me"`         // 自己发送的为1，联系人发送的为0
	FromLang             string `gorm:"column:from_lang;type:varchar(10);default:null" json:"from_lang"`    // 用于自动翻译的原始语言
	ToLang               string `gorm:"column:to_lang;type:varchar(10);default:null" json:"to_lang"`        // 用于自动翻译的目标语言
	Type                 string `gorm:"column:type;type:varchar(50);default:null" json:"type"`              // 消息类型（如"event"、"image"、"text"）
	AesKey               string `gorm:"column:aes_key;type:varchar(255);default:null" json:"aes_key"`       // 暂时留空
	Content              string `gorm:"column:content;type:text" json:"content"`                            // 消息内容
	ContentTranslate     string `gorm:"column:content_translate;type:text" json:"content_translate"`        // 翻译后的消息内容
	SeatID               int64  `gorm:"column:seat_id;type:bigint;not null;default:0" json:"-"`             // 座席id
	MessageID            string `gorm:"column:message_id;type:varchar(50);default:null" json:"-"`           // 来自ws的messageID
	FromJID              string `gorm:"column:from_jid;type:varchar(50);not null;default:''" json:"-"`      // 来自ws的fromJID消息发送者
	UserJID              string `gorm:"column:user_jid;type:varchar(50);not null;default:''" json:"-"`      // 来自ws的UserJID消息接收者
	FromPlatform         string `json:"from_platform,omitempty" json:"-"`                                   //可选。来自哪个平台客户端（用户判断是否要转发等）
	CreateTime           int64  `gorm:"column:create_time;type:datetime;default:null" json:"-"`             // 创建时间
	UpdateTime           int64  `gorm:"column:update_time;type:datetime;default:null" json:"-"`             // 更新时间
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
