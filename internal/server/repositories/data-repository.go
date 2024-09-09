package repositories

import (
	"time"

	"github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	"github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/entities"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type DataRepository struct {
	db *gorm.DB
}

func NewDataRepository(db *gorm.DB) *DataRepository {
	return &DataRepository{
		db: db,
	}
}

func (r *DataRepository) List(request requests.DataList) ([]*models.DataInfo, error) {
	var dataInfos []*models.DataInfo

	query := r.db.Table((&entities.Data{}).TableName()).Select(`
		   datas.id          as id,
		   datas.type        as type,
		   datas.value       as value,
		   datas.description as description`).
		Where("datas.deleted_at IS NULL").
		Order("id ASC")

	if request.Type != 0 {
		query = query.Where("datas.type=?", request.Type)
	}

	if request.UserID != 0 {
		query = query.Where("datas.user_id=?", request.UserID)
	}

	err := query.Find(&dataInfos).Error

	return dataInfos, err
}

func (r *DataRepository) Create(dataCreate requests.DataModel) (*models.DataInfo, error) {
	data := &entities.Data{
		UserID:      dataCreate.UserID,
		Type:        dataCreate.Type,
		Description: dataCreate.Description,
		Value:       dataCreate.Value,
	}

	err := r.db.
		Model(data).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"deleted_at", "value", "description"}),
		}).
		Create(data).Error
	if err != nil {
		return nil, err
	}

	return &models.DataInfo{
		ID:          data.ID,
		Type:        data.Type,
		Description: data.Description,
		Value:       data.Value,
	}, nil
}

func (r *DataRepository) Find(id uint, userID uint) (*models.DataInfo, error) {
	data := &models.DataInfo{}

	if err := r.db.
		Table((&entities.Data{}).TableName()).
		Select(`id,
  					  type,
                      description,
					  value`).
		Where("id = ?", id).
		Where("deleted_at is null").
		Where("user_id = ?", userID).
		Scan(&data).
		Error; err != nil {
		return nil, err
	}

	if data.ID == 0 {
		return nil, &NotFoundError{
			err: errNotFound,
		}
	}

	return data, nil
}

func (r *DataRepository) Update(id uint, request requests.DataModel) (*models.DataInfo, error) {
	newValues := map[string]interface{}{
		"updated_at":  time.Now(),
		"deleted_at":  nil,
		"user_id":     request.UserID,
		"type":        request.Type,
		"description": request.Description,
		"value":       request.Value,
	}

	data := &entities.Data{}
	result := r.db.Model(data).
		Where("id = ?", id).
		Where("user_id = ?", request.UserID).
		Clauses(clause.Returning{}).
		Updates(newValues)

	return &models.DataInfo{
		ID:          data.ID,
		Type:        data.Type,
		Description: data.Description,
		Value:       data.Value,
	}, result.Error
}

func (r *DataRepository) Delete(id uint, userId uint) error {
	result := *r.db.
		Where("id = ?", id).
		Where("user_id = ?", userId).
		Delete(&entities.Data{})

	return result.Error
}
