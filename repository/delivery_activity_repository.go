package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type DeliveryActivityRepository interface {
	CreateActivity(tx *gorm.DB, deliveryID uint, description string) (*model.DeliveryActivity, error)
}

type deliveryActivityRepository struct{}

func NewDeliveryActivityRepository() DeliveryActivityRepository {
	return &deliveryActivityRepository{}
}

func (d *deliveryActivityRepository) CreateActivity(tx *gorm.DB, deliveryID uint, description string) (*model.DeliveryActivity, error) {
	var deliveryAct = &model.DeliveryActivity{}
	deliveryAct.Description = description
	deliveryAct.DeliveryID = deliveryID
	result := tx.Model(&deliveryAct).Clauses(clause.Returning{}).Create(&deliveryAct)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot create delivery activity")
	}
	return deliveryAct, nil
}
