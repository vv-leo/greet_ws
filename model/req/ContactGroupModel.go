package req

type GroupIdModel struct {
	GroupID string `json:"group_id" binding:"required"`
}

type GroupIdNameModel struct {
	GroupID string `json:"group_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

type GroupIdContactIdModel struct {
	GroupID   string `json:"group_id"`
	ContactID string `json:"contact_id" binding:"required"`
}
