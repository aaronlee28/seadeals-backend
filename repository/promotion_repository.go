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
}

type promotionRepository struct{}

func NewPromotionRepository() PromotionRepository {
	return &promotionRepository{}
}

func (p *promotionRepository) GetPromotionBySellerID(tx *gorm.DB, sellerID uint) ([]*model.Promotion, error) {
	var promotion []*model.Promotion
	result := tx.Where("seller_id = ?", sellerID).Find(&promotion)
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
	}
	result := tx.Create(&promotion)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Failed to create promotion")
	}
	return promotion, nil
}
