package repositories

import (
	"errors"

	"github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"

	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/data"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/responses"

	"github.com/ShukinDmitriy/GophKeeper/internal/server/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(userRegister requests.UserRegister) (*responses.UserInfo, error) {
	passwordHash, err := r.GeneratePasswordHash(userRegister.Password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Login:    userRegister.Login,
		Password: string(passwordHash),
	}

	tx := r.db.Begin()

	query := tx.Model(&entities.User{}).
		Create(&user)
	err = query.Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &responses.UserInfo{
		ID:    user.ID,
		Login: user.Login,
	}, nil
}

func (r *UserRepository) Find(id uint) (*responses.UserInfo, error) {
	userModel := &responses.UserInfo{}
	if err := r.db.
		Select(`
		    users.id                                as id,
		    users.login                             as login`).
		Table("users").
		Where("users.id = ?", id).
		Where("users.deleted_at is null").
		Scan(&userModel).
		Error; err != nil {
		return nil, err
	}

	return userModel, nil
}

func (r *UserRepository) FindBy(filter data.UserSearch) (*responses.UserInfo, error) {
	user := &responses.UserInfo{}

	query := r.db

	if filter.Login != "" {
		query = query.Where("\"users\".\"login\" = ?", filter.Login)
	}

	if err := query.
		Select(`
		    users.id       as id,
		    users.login    as login,
		    users.password as password`).
		Table("users").
		Where("users.deleted_at is null").
		Limit(1).
		Scan(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	if user.ID == 0 {
		return nil, nil
	}

	return user, nil
}

func (r *UserRepository) GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 8)
}
