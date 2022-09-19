package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProvinceService interface {
	GetProvinces() ([]*model.Province, error)
}

type provinceService struct {
	db                 *gorm.DB
	provinceRepository repository.ProvinceRepository
}

type ProvinceServiceConfig struct {
	DB                 *gorm.DB
	ProvinceRepository repository.ProvinceRepository
}

func NewProvinceService(config *ProvinceServiceConfig) ProvinceService {
	return &provinceService{
		db:                 config.DB,
		provinceRepository: config.ProvinceRepository,
	}
}

func (p *provinceService) GetProvinces() ([]*model.Province, error) {
	tx := p.db.Begin()
	provinces, err := p.provinceRepository.GetProvinces(tx)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return provinces, nil
}
