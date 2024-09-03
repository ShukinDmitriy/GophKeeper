package repositories

import (
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
)

type DataRepositoryInterface interface {
	List(request requests.DataList) ([]*models.DataInfo, error)
	Create(dataCreate requests.DataModel) (*models.DataInfo, error)
	Find(id uint, userID uint) (*models.DataInfo, error)
	Update(id uint, request requests.DataModel) (*models.DataInfo, error)
	Delete(id uint, userID uint) error
}
