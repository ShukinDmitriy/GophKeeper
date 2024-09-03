package requests

import "github.com/ShukinDmitriy/GophKeeper/internal/common/models"

type DataList struct {
	Type   models.DataType `json:"type" query:"type"`
	UserID uint            `json:"user_id" query:"user_id"`
}
