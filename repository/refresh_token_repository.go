package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(tx *gorm.DB, userID uint, token string) error
}

type refreshTokenRepository struct {
}

type RefreshTokenRepositoryConfig struct {
}

func NewRefreshTokenRepo(c *RefreshTokenRepositoryConfig) RefreshTokenRepository {
	return &refreshTokenRepository{}
}

func (b *refreshTokenRepository) CreateRefreshToken(tx *gorm.DB, userID uint, token string) error {
	var tokenRefresh model.RefreshToken
	result := tx.Model(&tokenRefresh).Where("user_id = ?", userID).First(&tokenRefresh)
	if result.Error == nil {
		tokenRefresh.Token = token
		result = tx.Model(&tokenRefresh).Updates(&tokenRefresh)
		return result.Error
	}
	tokenRefresh.Token = token
	tokenRefresh.UserID = userID
	result = tx.Create(&tokenRefresh)
	return result.Error
}
