package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
)

type SellerRepository interface {
	FindSellerByID(tx *gorm.DB, id uint) (*model.Seller, error)
}

type sellerRepository struct{}

func NewSellerRepository() SellerRepository {
	return &sellerRepository{}
}

func (r *sellerRepository) FindSellerByID(tx *gorm.DB, id uint) (*model.Seller, error) {
	var seller *model.Seller
	result := tx.Preload("Address").Preload("User").First(&seller, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return seller, nil
}
