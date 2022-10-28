package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
	"time"
)

type PromotionService interface {
	GetPromotionByUserID(id uint) ([]*dto.GetPromotionRes, error)
	CreatePromotion(id uint, req *dto.CreatePromotionArrayReq) ([]*dto.CreatePromotionRes, error)
	ViewDetailPromotionByID(id uint) (*dto.GetPromotionRes, error)
	UpdatePromotion(req *dto.PatchPromotionArrayReq, userID uint) ([]*dto.PatchPromotionRes, error)
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

func (p *promotionService) CreatePromotion(id uint, req *dto.CreatePromotionArrayReq) ([]*dto.CreatePromotionRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	seller, err := p.sellerRepo.FindSellerByUserID(tx, id)
	if err != nil {
		return nil, err
	}

	sellerID := seller.ID
	var retArray []*dto.CreatePromotionRes

	for _, promotionReq := range req.CreatePromotion {
		if promotionReq.AmountType == "percentage" && promotionReq.Amount > 100 {
			err = apperror.BadRequestError("percentage amount exceeds 100%")
			return nil, apperror.BadRequestError("percentage amount exceeds 100%")
		}
		if !(promotionReq.AmountType == "percentage" || promotionReq.AmountType == "nominal") {
			promotionReq.AmountType = "nominal"
		}

		if promotionReq.StartDate.Before(time.Now()) {
			return nil, apperror.BadRequestError("date before today")
		}
		if promotionReq.EndDate.Before(promotionReq.StartDate) {
			return nil, apperror.BadRequestError("start date before end date")
		}

		var product *model.Product
		product, err = p.productRepo.GetProductDetail(tx, promotionReq.ProductID)
		if err != nil {
			return nil, err
		}
		if product.SellerID != sellerID {
			err = apperror.BadRequestError("Tidak bisa membuat promosi produk seller lain")
			return nil, err
		}

		var createPromo *model.Promotion
		createPromo, err = p.promotionRepository.CreatePromotion(tx, &promotionReq, sellerID)
		if err != nil {
			return nil, err
		}

		ret := dto.CreatePromotionRes{
			ID:          createPromo.ID,
			ProductID:   createPromo.ProductID,
			Name:        createPromo.Name,
			Description: createPromo.Description,
			StartDate:   createPromo.StartDate,
			EndDate:     createPromo.EndDate,
			Quota:       createPromo.Quota,
			MaxOrder:    createPromo.MaxOrder,
			AmountType:  createPromo.AmountType,
			Amount:      createPromo.Amount,
			BannerURL:   createPromo.BannerURL,
		}
		var userArray []*model.SocialGraph
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
		retArray = append(retArray, &ret)
	}

	return retArray, nil
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

func (p *promotionService) UpdatePromotion(req *dto.PatchPromotionArrayReq, userID uint) ([]*dto.PatchPromotionRes, error) {
	tx := p.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	var reqArray []*dto.PatchPromotionRes
	for _, promotion := range req.PatchPromotion {
		var promo *model.Promotion
		promo, err = p.promotionRepository.ViewDetailPromotionByID(tx, promotion.PromotionID)
		if promo.Product.Seller.UserID != userID {
			err = apperror.UnauthorizedError("tidak bisa mengupdate promotion seller lain")
			return nil, err
		}

		updatePromotion := &model.Promotion{
			Name:        promotion.Name,
			Description: promotion.Description,
			StartDate:   promotion.StartDate,
			EndDate:     promotion.EndDate,
			Quota:       promotion.Quota,
			MaxOrder:    promotion.MaxOrder,
			AmountType:  promotion.AmountType,
			Amount:      promotion.Amount,
		}
		var updatedPromotion *model.Promotion
		updatedPromotion, err = p.promotionRepository.UpdatePromotion(tx, promotion.PromotionID, updatePromotion)
		if err != nil {
			return nil, err
		}
		updatePromoRes := new(dto.PatchPromotionRes).PatchFromPromotion(updatedPromotion)
		reqArray = append(reqArray, updatePromoRes)
	}

	return reqArray, nil
}
