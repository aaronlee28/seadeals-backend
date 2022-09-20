package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type DistrictRepository interface {
	GetDistrictsByCityID(*gorm.DB, uint) ([]*model.District, error)
}

type districtRepository struct{}

func NewDistrictRepository() DistrictRepository {
	return &districtRepository{}
}

func (s *districtRepository) GetDistrictsByCityID(tx *gorm.DB, cityID uint) ([]*model.District, error) {
	var districts []*model.District
	result := tx.Where("city_id = ?", cityID).Find(&districts)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch districts")
	}

	return districts, result.Error
}
