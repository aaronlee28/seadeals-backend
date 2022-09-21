package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type SocialGraphRepository interface {
	GetFollowerCountBySellerID(tx *gorm.DB, sellerID uint) (int64, error)
	GetFollowingCountByUserID(tx *gorm.DB, userID uint) (int64, error)
}

type socialGraphRepository struct{}

func NewSocialGraphRepository() SocialGraphRepository {
	return &socialGraphRepository{}
}

func (s *socialGraphRepository) GetFollowerCountBySellerID(tx *gorm.DB, sellerID uint) (int64, error) {
	var count int64
	result := tx.Model(&model.SocialGraph{}).Where("seller_id = ?", sellerID).Count(&count)
	if result.Error != nil {
		return 0, apperror.InternalServerError("Cannot get review count")
	}

	return count, nil
}

func (s *socialGraphRepository) GetFollowingCountByUserID(tx *gorm.DB, userID uint) (int64, error) {
	var count int64
	result := tx.Model(&model.SocialGraph{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		return 0, apperror.InternalServerError("Cannot get review count")
	}

	return count, nil
}
