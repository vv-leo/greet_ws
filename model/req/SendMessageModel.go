package req

// 发送消息
type SendMessageModel struct {
	SendId    string `json:"sendId" binding:"required"`    //消息唯一id (发消息时由前端生成)
	ContactId string `json:"ContactId" binding:"required"` //会话id
	Type      string `json:"type" binding:"required"`      //消息类型（如"event"、"image"、"text"）
	Content   string `json:"content" binding:"required"`   //消息内容
}
