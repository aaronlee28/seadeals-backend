package service

import (
	"gorm.io/gorm"
	"seadeals-backend/dto"
	"seadeals-backend/repository"
)

type PromotionService interface {
	GetPromotionByUserID(id uint) (*[]dto.GetPromotionRes, error)
}

type promotionService struct {
	db                  *gorm.DB
	promotionRepository repository.PromotionRepository
	sellerRepo          repository.SellerRepository
	productRepo         repository.ProductRepository
}

type PromotionServiceConfig struct {
	DB                  *gorm.DB
	PromotionRepository repository.PromotionRepository
	SellerRepo          repository.SellerRepository
	ProductRepo         repository.ProductRepository
}

func NewPromotionService(c *PromotionServiceConfig) PromotionService {
	return &promotionService{
		db:                  c.DB,
		promotionRepository: c.PromotionRepository,
		sellerRepo:          c.SellerRepo,
		productRepo:         c.ProductRepo,
	}
}

func (p *promotionService) GetPromotionByUserID(id uint) (*[]dto.GetPromotionRes, error) {
	tx := p.db.Begin()
	seller, err := p.sellerRepo.FindSellerByUserID(tx, id)
	if err != nil {
		return nil, err
	}
	sellerID := seller.ID
	prs, err2 := p.promotionRepository.GetPromotionBySellerID(tx, sellerID)
	if err2 != nil {
		return nil, err2
	}
	var promoRes []dto.GetPromotionRes
	for _, promotion := range *prs {
		pr := new(dto.GetPromotionRes).FromPromotion(&promotion)
		photo, err3 := p.productRepo.GetProductPhotoURL(tx, promotion.ProductID)
		if err3 != nil {
			return nil, err3
		}
		pr.ProductPhotoURL = photo
		promoRes = append(promoRes, *pr)
	}
	return &promoRes, nil
}
