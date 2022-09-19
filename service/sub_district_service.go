package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type SubDistrictService interface {
	GetSubDistrictsByDistrictID(uint) ([]*model.SubDistrict, error)
}

type subDistrictService struct {
	db                    *gorm.DB
	subDistrictRepository repository.SubDistrictRepository
}

type SubDistrictServiceConfig struct {
	DB                    *gorm.DB
	SubDistrictRepository repository.SubDistrictRepository
}

func NewSubDistrictService(config *SubDistrictServiceConfig) SubDistrictService {
	return &subDistrictService{
		db:                    config.DB,
		subDistrictRepository: config.SubDistrictRepository,
	}
}

func (c *subDistrictService) GetSubDistrictsByDistrictID(districtID uint) ([]*model.SubDistrict, error) {
	tx := c.db.Begin()
	cities, err := c.subDistrictRepository.GetSubDistrictByDistrictID(tx, districtID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return cities, nil
}
