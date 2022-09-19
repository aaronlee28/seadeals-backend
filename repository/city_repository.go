package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type CityRepository interface {
	GetCitiesByProvinceID(*gorm.DB, uint) ([]*model.City, error)
}

type cityRepository struct {
}

type CityRepositoryConfig struct {
}

func NewCityRepository(c *CityRepositoryConfig) CityRepository {
	return &cityRepository{}
}

func (c *cityRepository) GetCitiesByProvinceID(tx *gorm.DB, provinceID uint) ([]*model.City, error) {
	var cities []*model.City
	result := tx.Where("province_id = ?", provinceID).Find(&cities)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch cities")
	}

	return cities, result.Error
}
