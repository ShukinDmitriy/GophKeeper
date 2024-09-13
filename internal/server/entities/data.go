package entities

import (
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	UserID      uint            `gorm:"type:bigint;not null"`
	Users       User            `gorm:"foreignKey:UserID;references:ID"`
	Type        models.DataType `json:"type" gorm:"type:integer;not null"`
	Value       string          `json:"value" gorm:"type:varchar;not null"`
	Description string          `json:"description" gorm:"type:varchar"`
}

func (d *Data) TableName() string {
	return "datas"
}
