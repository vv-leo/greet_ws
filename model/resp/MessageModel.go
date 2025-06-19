package resp

import (
	"git.exclouds.org/ww/global/model"
)

// MessageModel 发送消息-响应
/*type MessageModel struct {
	Ack              int       `json:"ack"`               // 标识是否已读的标志，暂时固定为1
	WaAccountID      int64     `json:"wa_account_id"`     // 发送消息的WS账号在系统中的id编号
	To               string  `json:"to"`                // 对话id
	ID               string    `json:"id"`                // 消息的唯一id，可以是数字或字符串
	Status           string    `json:"status"`            // 消息状态 string | number
	ToContactID      string  `json:"toContactId"`       // 同对话id
	SendTime         int64     `json:"sendTime"`          // 发送时间的时间戳
	FromUser         UserModel `json:"fromUser"`          // 发送消息的用户信息
	IsFromMe         int       `json:"is_from_me"`        // 自己发送的为1，联系人发送的为0
	FromLang         string    `json:"from_lang"`         // 用于自动翻译的原始语言
	ToLang           string    `json:"to_lang"`           // 用于自动翻译的目标语言
	Type             string    `json:"type"`              // 消息类型（如"event"、"image"、"text"）
	AesKey           string    `json:"aes_key"`           // 暂时留空
	Content          string    `json:"content"`           // 消息内容
	ContentTranslate string    `json:"content_translate"` // 翻译后的消息内容
}*/

type UserModel struct {
	ID          string `json:"id"`          // 如果是自己发送的这里是1，如果是联系人发送的则显示对话id
	DisplayName string `json:"displayName"` // 昵称/备注
	Avatar      string `json:"avatar"`      // 头像地址，没有则留空
}

type MessageWithUserModel struct {
	model.MessageModel
	FromUser UserModel `json:"fromUser"`
}

type RelayMessageModel struct {
	Contacts []ContactModel         `json:"contacts"`
	Messages []MessageWithUserModel `json:"messages"`
}
