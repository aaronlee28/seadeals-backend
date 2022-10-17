package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type CourierRepository interface {
	GetAllCouriers(tx *gorm.DB) ([]*model.Courier, error)
	GetCourierDetailByID(tx *gorm.DB, courierID uint) (*model.Courier, error)
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

func (c *courierRepository) GetCourierDetailByID(tx *gorm.DB, courierID uint) (*model.Courier, error) {
	var courier = &model.Courier{}
	courier.ID = courierID
	result := tx.Model(&courier).Find(&courier)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.BadRequestError("No such courier exists")
		}
		return nil, apperror.InternalServerError("Cannot fetch courier")
	}
	return courier, nil
}
