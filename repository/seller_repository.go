package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type SellerRepository interface {
	FindSellerDetailByID(tx *gorm.DB, sellerID uint, userID uint) (*model.Seller, error)
	FindSellerByID(tx *gorm.DB, sellerID uint) (*model.Seller, error)
	FindSellerByUserID(tx *gorm.DB, userID uint) (*model.Seller, error)
}

type sellerRepository struct{}

func NewSellerRepository() SellerRepository {
	return &sellerRepository{}
}

func (r *sellerRepository) FindSellerDetailByID(tx *gorm.DB, sellerID uint, userID uint) (*model.Seller, error) {
	var seller *model.Seller
	result := tx.Preload("Address").Preload("User").Preload("SocialGraph", "seller_id = ? AND user_id = ? AND is_follow IS TRUE", sellerID, userID).First(&seller, sellerID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, apperror.NotFoundError("No such seller exists")
		}
		return nil, apperror.InternalServerError("Cannot fetch seller detail")
	}
	return seller, nil
}

func (r *sellerRepository) FindSellerByID(tx *gorm.DB, sellerID uint) (*model.Seller, error) {
	var seller *model.Seller
	result := tx.Preload("Address").Preload("User").First(&seller, sellerID)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError("No such seller exists")
	}
	return seller, result.Error
}

func (r *sellerRepository) FindSellerByUserID(tx *gorm.DB, userID uint) (*model.Seller, error) {
	var seller *model.Seller
	result := tx.Where("user_id = ?", userID).Preload("Address").Preload("User").First(&seller)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, apperror.NotFoundError("Seller doesn't exists")
	}
	return seller, result.Error
}
