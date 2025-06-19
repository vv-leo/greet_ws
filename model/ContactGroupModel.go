package model

// ContactGroupModel 表示联系人分组表
type ContactGroupModel struct {
	Id         string `json:"id" db:"id"`                  // 分组ID
	SeatID     int64  `json:"-" db:"seat_id"`              // 座席ID
	Name       string `json:"name" db:"name"`              // 分组名
	Count      int    `json:"count" db:"count"`            // 分组名
	CreateTime int64  `json:"created_at" db:"create_time"` // 创建时间
	UpdateTime int64  `json:"updated_at" db:"update_time"` // 修改时间
}

func (c ContactGroupModel) TableName() string {
	return "contact_group"
}
