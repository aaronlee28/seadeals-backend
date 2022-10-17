package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type FavoriteRepository interface {
	FavoriteToProduct(tx *gorm.DB, userID uint, productID uint) (*model.Favorite, error)
	GetUserFavoriteCount(tx *gorm.DB, userID uint) (uint, error)
}

type favoriteRepository struct{}

func NewFavoriteRepository() FavoriteRepository {
	return &favoriteRepository{}
}

func (f *favoriteRepository) FavoriteToProduct(tx *gorm.DB, userID uint, productID uint) (*model.Favorite, error) {
	var favorite = &model.Favorite{}
	result := tx.Model(&favorite).Where("user_id = ?", userID).Where("product_id = ?", productID).First(&favorite)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return nil, apperror.InternalServerError("Cannot find favorite status to a product")
		}
		favorite.UserID = userID
		favorite.ProductID = productID
		favorite.IsFavorite = true
		result = tx.Create(favorite)
		if result.Error != nil {
			return nil, apperror.InternalServerError("Cannot favorite a product")
		}
		return favorite, nil
	}

	currentFavorite := !favorite.IsFavorite
	result = tx.Model(&favorite).Clauses(clause.Returning{}).Updates(map[string]interface{}{"is_favorite": currentFavorite})
	if result.Error != nil {
		return nil, apperror.InternalServerError("Cannot favorite a product")
	}

	return favorite, nil
}

func (f *favoriteRepository) GetUserFavoriteCount(tx *gorm.DB, userID uint) (uint, error) {
	var countInt64 int64
	result := tx.Model(&model.Favorite{}).
		Where("is_favorite = true AND user_id = ?", userID).
		Count(&countInt64)

	if result.Error != nil {
		return 0, apperror.InternalServerError("Cannot get favorite count")
	}

	count := uint(countInt64)

	return count, nil
}
