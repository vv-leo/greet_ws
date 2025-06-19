package model

type WaAccountModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	WaAccount   string `gorm:"type:varchar(255);unique;not null" json:"wa_account"`
	ToWaAccount string `gorm:"type:varchar(255);unique;not null" json:"to_wa_account"`
	FromType    int    `gorm:"type:tinyint;default:0" json:"from_type"` // 0:打招呼；1:陌生人来联系
	SeatID      int64  `gorm:"type:tinyint;default:0" json:"seat_id"`   // 座席ID
	CreateTime  int64  `gorm:"type:datetime;default:null" json:"create_time"`
	UpdateTime  int64  `gorm:"type:datetime;default:null" json:"update_time"`
}

func (WaAccountModel) TableName() string {
	return "wa_account"
}
