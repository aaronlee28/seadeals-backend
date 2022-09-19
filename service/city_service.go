package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type CityService interface {
	GetCitiesByProvinceID(uint) ([]*model.City, error)
}

type cityService struct {
	db             *gorm.DB
	cityRepository repository.CityRepository
}

type CityServiceConfig struct {
	DB             *gorm.DB
	CityRepository repository.CityRepository
}

func NewCityService(config *CityServiceConfig) CityService {
	return &cityService{
		db:             config.DB,
		cityRepository: config.CityRepository,
	}
}

func (c *cityService) GetCitiesByProvinceID(provinceID uint) ([]*model.City, error) {
	tx := c.db.Begin()
	cities, err := c.cityRepository.GetCitiesByProvinceID(tx, provinceID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return cities, nil
}
