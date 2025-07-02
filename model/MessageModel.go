package model

type MessageModel struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	ContactId   string `json:"contactId"`
	SendTime    int64  `json:"sendTime"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
	IsFromMe    int8   `json:"isFromMe"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	SeatId      int64  `json:"seatId"`
	MessageId   string `json:"messageId"`
	TargetAcc   string `json:"targetAcc"`
	PlatformAcc string `json:"platformAcc"`
	CreateTime  int64  `json:"createTime"`
	UpdateTime  int64  `json:"updateTime"`
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
