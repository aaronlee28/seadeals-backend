package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type CourierRepository interface {
	GetAllCouriers(tx *gorm.DB) ([]*model.Courier, error)
}

type courierRepository struct{}

func NewCourierRepository() CourierRepository {
	return &courierRepository{}
}

func (c *courierRepository) GetAllCouriers(tx *gorm.DB) ([]*model.Courier, error) {
	var couriers []*model.Courier
	result := tx.Model(&couriers).Find(&couriers)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot fetch couriers")
	}
	return couriers, nil
}
