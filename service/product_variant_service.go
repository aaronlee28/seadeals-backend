package service

import (
	"errors"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type ProductVariantService interface {
	FindAllProductVariantByProductID(productID uint) (*dto.ProductVariantRes, error)
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

func (s *productVariantService) FindAllProductVariantByProductID(productID uint) (*dto.ProductVariantRes, error) {
	tx := s.db.Begin()

	productVariants, err := s.productVariantRepo.FindAllProductVariantByProductID(tx, productID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, &apperror.ProductNotFoundError{}) {
			return nil, apperror.NotFoundError(err.Error())
		}
		return nil, err
	}

	var productVariantRes []*dto.GetProductVariantRes
	minPrice := productVariants[0].Price
	maxPrice := minPrice

	for _, pv := range productVariants {
		if pv.Price > maxPrice {
			maxPrice = pv.Price
		}
		if pv.Price < minPrice {
			minPrice = pv.Price
		}
		productVariantRes = append(productVariantRes, new(dto.GetProductVariantRes).From(pv))
	}

	res := &dto.ProductVariantRes{
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
		ProductVariants: productVariantRes,
	}

	tx.Commit()
	return res, nil
}
