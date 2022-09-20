package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type SubDistrictRepository interface {
	GetSubDistrictByDistrictID(*gorm.DB, uint) ([]*model.SubDistrict, error)
}

type subDistrictRepository struct{}

func NewSubDistrictRepository() SubDistrictRepository {
	return &subDistrictRepository{}
}

func (s *subDistrictRepository) GetSubDistrictByDistrictID(tx *gorm.DB, districtID uint) ([]*model.SubDistrict, error) {
	var subDistricts []*model.SubDistrict
	result := tx.Where("district_id = ?", districtID).Find(&subDistricts)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch sub districts")
	}

	return subDistricts, result.Error
}
