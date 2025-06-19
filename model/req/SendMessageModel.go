package req

// 发送消息
type SendMessageModel struct {
	UUID    string `json:"uuid" binding:"required"`    //消息唯一id (发消息时由前端生成)
	To      string `json:"to"`                         //对话id
	Type    string `json:"type" binding:"required"`    //消息类型（如"event"、"image"、"text"）
	Content string `json:"content" binding:"required"` //消息内容
}

type LastMessageTimeContactIdModel struct {
	LastMessageTime string `json:"last_message_time"` //最新消息时间
	ContactId       string `json:"contact_id"`        //对话id
}

type LastMessagePageContactIdModel struct {
	ContactId   string `json:"contact_id"`  // 对话ID
	ToContactId string `json:"toContactId"` // 对话ID（目标联系人ID）
	Page        int    `json:"page"`        // 分页参数：当前页码
	Limit       int    `json:"limit"`       // 分页参数：每页最大记录数
	IsGroup     int    `json:"is_group"`    //
	LastId      string `json:"last_id"`     // 最后一条消息的ID
}
