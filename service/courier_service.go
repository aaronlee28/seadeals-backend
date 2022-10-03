package service

import (
	"gorm.io/gorm"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type CourierService interface {
	GetAllCouriers() ([]*model.Courier, error)
}

type courierService struct {
	db                *gorm.DB
	courierRepository repository.CourierRepository
}

type CourierServiceConfig struct {
	DB                *gorm.DB
	CourierRepository repository.CourierRepository
}

func NewCourierService(c *CourierServiceConfig) CourierService {
	return &courierService{
		db:                c.DB,
		courierRepository: c.CourierRepository,
	}
}

func (c *courierService) GetAllCouriers() ([]*model.Courier, error) {
	tx := c.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	couriers, err := c.courierRepository.GetAllCouriers(tx)
	if err != nil {
		return nil, err
	}

	return couriers, nil
}
