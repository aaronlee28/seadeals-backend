package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
)

type PromotionRepository interface {
}

type promotionRepository struct {
}

func NewPromotionRepository() PromotionRepository {
	return &promotionRepository{}
}

func (r *productRepository) GetPromotionByID(tx *gorm.DB, id uint) (*[]model.Promotion, error) {
	var promotion *[]model.Promotion
	result := tx.Where("id = ?").Find(&promotion)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
