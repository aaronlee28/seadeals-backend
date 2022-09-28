package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
)

type PromotionRepository interface {
	GetPromotionBySellerID(tx *gorm.DB, sellerID uint) (*[]model.Promotion, error)
}

type promotionRepository struct{}

func NewPromotionRepository() PromotionRepository {
	return &promotionRepository{}
}

func (p *promotionRepository) GetPromotionBySellerID(tx *gorm.DB, sellerID uint) (*[]model.Promotion, error) {
	var promotion *[]model.Promotion
	result := tx.Where("seller_id = ?", sellerID).Find(&promotion)
	if result.Error != nil {
		return nil, result.Error
	}
	return promotion, nil
}