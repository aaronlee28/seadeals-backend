package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type ProductVariantService interface {
	FindAllProductVariantByProductID(productID uint) (*dto.ProductVariantRes, error)
	GetVariantPriceAfterPromotionByProductID(productID int) (*dto.ProductVariantPriceRes, error)
}

type productVariantService struct {
	db                 *gorm.DB
	productRepo        repository.ProductRepository
	productVariantRepo repository.ProductVariantRepository
	productVarDetRepo  repository.ProductVariantDetailRepository
}

type ProductVariantServiceConfig struct {
	DB                 *gorm.DB
	ProductRepo        repository.ProductRepository
	ProductVariantRepo repository.ProductVariantRepository
	ProductVarDetRepo  repository.ProductVariantDetailRepository
}

func NewProductVariantService(c *ProductVariantServiceConfig) ProductVariantService {
	return &productVariantService{
		db:                 c.DB,
		productRepo:        c.ProductRepo,
		productVariantRepo: c.ProductVariantRepo,
		productVarDetRepo:  c.ProductVarDetRepo,
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

func (s *productVariantService) GetVariantPriceAfterPromotionByProductID(productID int) (*dto.ProductVariantPriceRes, error) {
	tx := s.db.Begin()
	id := uint(productID)

	product, err := s.productRepo.GetProductDetail(tx, id)
	fmt.Println("prodssss", product.Name)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var variants []*dto.ProductVariantPromotionRes
	for _, variant := range product.ProductVariantDetail {
		vr := new(dto.ProductVariantPromotionRes).FromProductVariantDetail(variant)
		vr.PriceAfterPromotion = vr.Price - product.Promotion.Amount
		variants = append(variants, vr)
	}
	res := dto.ProductVariantPriceRes{
		ProductID:        product.ID,
		ProductName:      product.Name,
		ProductPromotion: product.Promotion.Amount,
		ProductVariant:   variants,
	}

	return &res, nil

}
