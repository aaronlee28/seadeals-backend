package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ProvinceRepository interface {
	GetProvinces(*gorm.DB) ([]*model.Province, error)
}

type provinceRepository struct {
}

type ProvinceRepositoryConfig struct {
}

func NewProvinceRepository(c *ProvinceRepositoryConfig) ProvinceRepository {
	return &provinceRepository{}
}

func (s *provinceRepository) GetProvinces(tx *gorm.DB) ([]*model.Province, error) {
	var provinces []*model.Province
	result := tx.Find(&provinces)
	if result.Error != nil {
		return nil, apperror.InternalServerError("cannot fetch provinces")
	}

	return provinces, result.Error
}
