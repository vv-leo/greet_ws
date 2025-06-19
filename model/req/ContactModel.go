package req

type IdRemarkModel struct {
	Id     string `binding:"required"`
	Remark string `binding:"omitempty"`
}

type GetContactsModel struct {
	StartTime  string `json:"start_time"`  // 查询的开始时间，格式通常为 YYYY-MM-DD HH:mm:ss
	Ver        int    `json:"ver"`         // 版本号，例如 ver: 2
	ReadType   string `json:"read_type"`   // 读取状态: ""(全部), "is_read"(已读), "unread"(未读)
	Blocked    string `json:"blocked"`     // 拉黑状态: ""(未拉黑), "1"(拉黑)
	Phone      string `json:"phone"`       // 搜索手机号
	AllowReply string `json:"allow_reply"` // 回复状态: "all"(全部), "reply"(已回复), "unreply"(未回复)
	GroupID    string `json:"group_id"`    // 分组ID: ""(全部) 或具体分组 ID
}
