package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type DistrictService interface {
	GetDistrictsByCityID(uint) ([]*model.District, error)
}

type districtService struct {
	db                 *gorm.DB
	districtRepository repository.DistrictRepository
}

type DistrictServiceConfig struct {
	DB                 *gorm.DB
	DistrictRepository repository.DistrictRepository
}

func NewDistrictService(config *DistrictServiceConfig) DistrictService {
	return &districtService{
		db:                 config.DB,
		districtRepository: config.DistrictRepository,
	}
}

func (c *districtService) GetDistrictsByCityID(cityID uint) ([]*model.District, error) {
	tx := c.db.Begin()
	cities, err := c.districtRepository.GetDistrictsByCityID(tx, cityID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return cities, nil
}
