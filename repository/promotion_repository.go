package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/dto"
	"seadeals-backend/model"
)

type PromotionRepository interface {
	GetPromotionBySellerID(tx *gorm.DB, sellerID uint) ([]*model.Promotion, error)
	CreatePromotion(tx *gorm.DB, req *dto.CreatePromotionReq, sellerID uint) (*model.Promotion, error)
	ViewDetailPromotionByID(tx *gorm.DB, id uint) (*model.Promotion, error)
	UpdatePromotion(tx *gorm.DB, promoID uint, updatePromo *model.Promotion) (*model.Promotion, error)
}

type promotionRepository struct{}

func NewPromotionRepository() PromotionRepository {
	return &promotionRepository{}
}

func (p *promotionRepository) GetPromotionBySellerID(tx *gorm.DB, sellerID uint) ([]*model.Promotion, error) {
	var promotion []*model.Promotion
	result := tx.Model(&promotion).Where("seller_id = ?", sellerID).Preload("Product").Find(&promotion)
	return promotion, result.Error
}

func (p *promotionRepository) CreatePromotion(tx *gorm.DB, req *dto.CreatePromotionReq, sellerID uint) (*model.Promotion, error) {
	promotion := &model.Promotion{
		ProductID:   req.ProductID,
		SellerID:    sellerID,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Quota:       req.Quota,
		MaxOrder:    req.MaxOrder,
		AmountType:  req.AmountType,
		Amount:      req.Amount,
		BannerURL:   req.BannerURL,
	}
	result := tx.Create(&promotion)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Failed to create promotion")
	}
	return promotion, nil
}

func (p *promotionRepository) ViewDetailPromotionByID(tx *gorm.DB, id uint) (*model.Promotion, error) {
	var promotion *model.Promotion
	result := tx.Where("id = ?", id).Preload("Product.Seller").First(&promotion)
	return promotion, result.Error
}

func (p *promotionRepository) UpdatePromotion(tx *gorm.DB, promoID uint, updatePromo *model.Promotion) (*model.Promotion, error) {
	var updatedPromotion *model.Promotion
	result := tx.First(&updatedPromotion, promoID).Updates(&updatePromo)
	return updatedPromotion, result.Error
}
