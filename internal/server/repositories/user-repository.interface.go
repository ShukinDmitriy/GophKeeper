package repositories

import (
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/data"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/responses"
)

type UserRepositoryInterface interface {
	Create(userRegister requests.UserRegister) (*responses.UserInfo, error)
	Find(id uint) (*responses.UserInfo, error)
	FindBy(filter data.UserSearch) (*responses.UserInfo, error)
	GeneratePasswordHash(password string) ([]byte, error)
}
