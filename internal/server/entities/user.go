package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login    string `json:"login" gorm:"type:varchar;not null;unique"`
	Password string `json:"password" gorm:"type:varchar;not null"`
}
