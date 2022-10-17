package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type SocialGraphRepository interface {
	GetFollowerCountBySellerID(tx *gorm.DB, sellerID uint) (int64, error)
	GetFollowingCountByUserID(tx *gorm.DB, userID uint) (int64, error)
	FollowToSeller(tx *gorm.DB, userID uint, sellerID uint) (*model.SocialGraph, error)
	GetFavoriteUserID(tx *gorm.DB, productID uint) ([]*model.Favorite, error)
	GetFollowerUserID(tx *gorm.DB, sellerID uint) ([]*model.SocialGraph, error)
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

func (s *socialGraphRepository) FollowToSeller(tx *gorm.DB, userID uint, sellerID uint) (*model.SocialGraph, error) {
	var socialGraph = &model.SocialGraph{}
	result := tx.Model(&socialGraph).Where("user_id = ?", userID).Where("seller_id = ?", sellerID).First(&socialGraph)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, apperror.InternalServerError("Cannot find follow status on a seller")
		}
		socialGraph.UserID = userID
		socialGraph.SellerID = sellerID
		socialGraph.IsFollow = true
		result = tx.Create(socialGraph)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot social graph an item")
		}
	}

	currentFollow := !socialGraph.IsFollow
	result = tx.Model(&socialGraph).Clauses(clause.Returning{}).Updates(map[string]interface{}{"is_follow": currentFollow})
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot social graph an item")
	}

	return socialGraph, nil

}

func (s *socialGraphRepository) GetFavoriteUserID(tx *gorm.DB, productID uint) ([]*model.Favorite, error) {
	var userArray []*model.Favorite
	result := tx.Clauses(clause.Returning{}).Where("product_id = ?", productID).Where("is_favorite = true").Find(&userArray)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot get users")
	}

	return userArray, nil
}

func (s *socialGraphRepository) GetFollowerUserID(tx *gorm.DB, sellerID uint) ([]*model.SocialGraph, error) {
	var userArray []*model.SocialGraph
	result := tx.Clauses(clause.Returning{}).Where("seller_id = ?", sellerID).Where("is_follow = true").Find(&userArray)
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot get users")
	}

	return userArray, nil
}
