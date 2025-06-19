package model

// ContactGroupRelationModel 表示联系人分组关系表
type ContactGroupRelationModel struct {
	ID         uint64 `json:"id" db:"id"`                                               // 关系表ID
	GroupID    int64  `json:"group_id" db:"group_id"`                                   // 分组ID
	SeatID     int64  `json:"seat_id" db:"seat_id"`                                     // 座席ID
	ContactID  string `gorm:"column:contact_id;type:bigint;not null" json:"contact_id"` // 对话ID
	CreateTime int64  `json:"create_time" db:"create_time"`                             // 创建时间
	UpdateTime int64  `json:"update_time" db:"update_time"`                             // 修改时间
}

func (ContactGroupRelationModel) TableName() string {
	return "contact_group_relation"
}
