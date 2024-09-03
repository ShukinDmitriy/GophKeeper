package requests

import "github.com/ShukinDmitriy/GophKeeper/internal/common/models"

type DataModel struct {
	Type        models.DataType `json:"type"`
	Description string          `json:"description"`
	Value       string          `json:"value" validate:"required"`
	UserID      uint            `json:"user_id" validate:"required"`
}
