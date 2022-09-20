package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type ProductVariantService interface {
	FindAllProductVariantByProductID(productID uint) ([]*dto.GetProductVariantRes, error)
}

type productVariantService struct {
	db                 *gorm.DB
	productVariantRepo repository.ProductVariantRepository
}

type ProductVariantServiceConfig struct {
	DB                 *gorm.DB
	ProductVariantRepo repository.ProductVariantRepository
}

func NewProductVariantService(c *ProductVariantServiceConfig) ProductVariantService {
	return &productVariantService{
		db:                 c.DB,
		productVariantRepo: c.ProductVariantRepo,
	}
}

func (s *productVariantService) FindAllProductVariantByProductID(productID uint) ([]*dto.GetProductVariantRes, error) {
	tx := s.db.Begin()

	productVariants, err := s.productVariantRepo.FindAllProductVariantByProductID(tx, productID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var res []*dto.GetProductVariantRes
	for _, pv := range productVariants {
		res = append(res, new(dto.GetProductVariantRes).From(pv))
	}

	tx.Commit()
	return res, nil
}
