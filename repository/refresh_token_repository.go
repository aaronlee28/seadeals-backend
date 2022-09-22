package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(tx *gorm.DB, userID uint, token string) error
	DeleteRefreshToken(tx *gorm.DB, userID uint) error
	CheckIfTokenExist(tx *gorm.DB, token string) (bool, uint, error)
}

type refreshTokenRepository struct{}

func NewRefreshTokenRepo() RefreshTokenRepository {
	return &refreshTokenRepository{}
}

func (r *refreshTokenRepository) CreateRefreshToken(tx *gorm.DB, userID uint, token string) error {
	var tokenRefresh model.RefreshToken
	result := tx.Model(&tokenRefresh).Where("user_id = ?", userID).First(&tokenRefresh)
	if result.Error == nil {
		if result.Error != gorm.ErrRecordNotFound {
			tokenRefresh.Token = token
			result = tx.Model(&tokenRefresh).Updates(&tokenRefresh)
			return nil
		}
		return apperror.InternalServerError("Cannot find refresh token")
	}
	tokenRefresh.Token = token
	tokenRefresh.UserID = userID
	result = tx.Create(&tokenRefresh)
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(tx *gorm.DB, userID uint) error {
	var tokenRefresh = &model.RefreshToken{}
	result := tx.Model(&tokenRefresh).Where("user_id = ?", userID).First(&tokenRefresh)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return apperror.NotFoundError("Refresh token not found")
		}
		return apperror.InternalServerError("Cannot get token")
	}
	result = tx.Delete(&tokenRefresh)
	if result.Error != nil {
		return apperror.InternalServerError("Cannot delete Refresh Token")
	}
	return nil
}

func (r *refreshTokenRepository) CheckIfTokenExist(tx *gorm.DB, token string) (bool, uint, error) {
	var tokenRefresh = &model.RefreshToken{}
	result := tx.Model(&tokenRefresh).Where("token = ?", token).First(&tokenRefresh)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, 0, apperror.NotFoundError("Refresh token not found")
		}
		return false, 0, apperror.InternalServerError("Cannot get token")
	}
	return true, tokenRefresh.UserID, nil
}
