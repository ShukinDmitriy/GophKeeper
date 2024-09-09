package http

import (
	"context"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	commonRequests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
)

type ClientInterface interface {
	Login(ctx context.Context, data commonRequests.UserLogin) error
	Register(ctx context.Context, data commonRequests.UserRegister) error
	GetList(ctx context.Context, dataType models.DataType) ([]models.DataInfo, error)
	CreateData(ctx context.Context, data commonRequests.DataModel) (*models.DataInfo, error)
	UpdateData(ctx context.Context, data models.DataInfo) (*models.DataInfo, error)
	DeleteData(ctx context.Context, id uint) error
}
