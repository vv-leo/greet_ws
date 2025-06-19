package model

import (
	"gorm.io/gorm"
)

type JwtBlacklistModel struct {
	gorm.Model
	Jwt string `gorm:"type:text;comment:jwt"`
}
