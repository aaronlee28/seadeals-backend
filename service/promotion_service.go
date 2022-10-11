package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type PromotionService interface {
	GetPromotionByUserID(id uint) ([]*dto.GetPromotionRes, error)
	CreatePromotion(id uint, req *dto.CreatePromotionReq) (*dto.CreatePromotionRes, error)
	ViewDetailPromotionByID(id uint) (*dto.GetPromotionRes, error)
	UpdatePromotion(req *dto.PatchPromotionReq, promoID uint, userID uint) (*dto.PatchPromotionRes, error)
}

type promotionService struct {
	db                  *gorm.DB
	promotionRepository repository.PromotionRepository
	sellerRepo          repository.SellerRepository
	productRepo         repository.ProductRepository
	socialGraphRepo     repository.SocialGraphRepository
	notificationRepo    repository.NotificationRepository
}

type PromotionServiceConfig struct {
	DB                  *gorm.DB
	PromotionRepository repository.PromotionRepository
	SellerRepo          repository.SellerRepository
	ProductRepo         repository.ProductRepository
	SocialGraphRepo     repository.SocialGraphRepository
	NotificationRepo    repository.NotificationRepository
}

func NewPromotionService(c *PromotionServiceConfig) PromotionService {
	return &promotionService{
		db:                  c.DB,
		promotionRepository: c.PromotionRepository,
		sellerRepo:          c.SellerRepo,
		productRepo:         c.ProductRepo,
		socialGraphRepo:     c.SocialGraphRepo,
		notificationRepo:    c.NotificationRepo,
	}
}

func (p *promotionService) GetPromotionByUserID(id uint) ([]*dto.GetPromotionRes, error) {
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

	var promoRes = make([]*dto.GetPromotionRes, 0)
	for _, promotion := range prs {
		pr := new(dto.GetPromotionRes).FromPromotion(promotion)
		var photo string
		photo, err = p.productRepo.GetProductPhotoURL(tx, promotion.ProductID)
		if err != nil {
			return nil, err
		}
		pr.ProductPhotoURL = photo
		promoRes = append(promoRes, pr)
	}
	return promoRes, nil
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
		err = apperror.BadRequestError("percentage amount exceeds 100%")
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
	var userArray []*model.Favorite
	userArray, err = p.socialGraphRepo.GetFollowerUserID(tx, seller.ID)
	for _, user := range userArray {
		newNotification := &model.Notification{
			UserID:   user.UserID,
			SellerID: seller.ID,
			Title:    dto.NotificationFollowPromosi,
			Detail:   "Seller adds new promotion",
		}
		p.notificationRepo.AddToNotificationFromModel(tx, newNotification)
	}
	return &ret, nil
}

func (p *promotionService) ViewDetailPromotionByID(id uint) (*dto.GetPromotionRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	promo, err := p.promotionRepository.ViewDetailPromotionByID(tx, id)
	var photo string
	photo, err = p.productRepo.GetProductPhotoURL(tx, promo.ProductID)
	if err != nil {
		return nil, err
	}
	promoRes := new(dto.GetPromotionRes).FromPromotion(promo)
	promoRes.ProductPhotoURL = photo

	return promoRes, nil
}

func (p *promotionService) UpdatePromotion(req *dto.PatchPromotionReq, promoID uint, userID uint) (*dto.PatchPromotionRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	promo, err := p.promotionRepository.ViewDetailPromotionByID(tx, promoID)
	if promo.Product.Seller.UserID != userID {
		err = apperror.UnauthorizedError("cannot update other shop promotion")
		return nil, err
	}

	updatePromotion := &model.Promotion{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Quota:       req.Quota,
		MaxOrder:    req.MaxOrder,
		AmountType:  req.AmountType,
		Amount:      req.Amount,
	}
	var updatedPromotion *model.Promotion
	updatedPromotion, err = p.promotionRepository.UpdatePromotion(tx, promoID, updatePromotion)
	if err != nil {
		return nil, err
	}
	updatePromoRes := new(dto.PatchPromotionRes).PatchFromPromotion(updatedPromotion)
	return updatePromoRes, nil
}
