package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/repository"
)

type PromotionService interface {
	GetPromotionByUserID(id uint) (*[]dto.GetPromotionRes, error)
	CreatePromotion(id uint, req *dto.CreatePromotionReq) (*dto.CreatePromotionRes, error)
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
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := p.sellerRepo.FindSellerByUserID(tx, id)
	if err != nil {
		return nil, err
	}

	sellerID := seller.ID
	prs, err := p.promotionRepository.GetPromotionBySellerID(tx, sellerID)
	if err != nil {
		return nil, err
	}

	var promoRes []dto.GetPromotionRes
	for _, promotion := range *prs {
		pr := new(dto.GetPromotionRes).FromPromotion(&promotion)
		var photo string
		photo, err = p.productRepo.GetProductPhotoURL(tx, promotion.ProductID)
		if err != nil {
			return nil, err
		}
		pr.ProductPhotoURL = photo
		promoRes = append(promoRes, *pr)
	}
	return &promoRes, nil
}

func (p *promotionService) CreatePromotion(id uint, req *dto.CreatePromotionReq) (*dto.CreatePromotionRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := p.sellerRepo.FindSellerByUserID(tx, id)
	if err != nil {
		return nil, err
	}

	sellerID := seller.ID
	if req.AmountType == "percentage" && req.Amount > 100 {
		return nil, apperror.BadRequestError("percentage amount exceeds 100%")
	}
	if !(req.AmountType == "percentage" || req.AmountType == "nominal") {
		req.AmountType = "nominal"
	}

	createPromo, err := p.promotionRepository.CreatePromotion(tx, req, sellerID)
	if err != nil {
		return nil, err
	}

	ret := dto.CreatePromotionRes{
		ID:        createPromo.ID,
		ProductID: createPromo.ProductID,
		Name:      createPromo.Name,
	}
	return &ret, nil
}
