package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type SellerService interface {
	FindSellerByID(id uint) (*dto.GetSellerRes, error)
}

type sellerService struct {
	db         *gorm.DB
	sellerRepo repository.SellerRepository
}

type SellerServiceConfig struct {
	DB         *gorm.DB
	SellerRepo repository.SellerRepository
}

func NewSellerService(c *SellerServiceConfig) SellerService {
	return &sellerService{
		db:         c.DB,
		sellerRepo: c.SellerRepo,
	}
}

func (s *sellerService) FindSellerByID(id uint) (*dto.GetSellerRes, error) {
	tx := s.db.Begin()
	seller, err := s.sellerRepo.FindSellerByID(tx, id)
	if err != nil {
		return nil, err
	}

	res := new(dto.GetSellerRes).From(seller)

	tx.Commit()
	return res, nil
}
