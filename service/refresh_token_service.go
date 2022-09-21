package service

import (
	"gorm.io/gorm"
	"seadeals-backend/repository"
)

type RefreshTokenService interface {
	CheckIfTokenExist(token string) (bool, uint, error)
}

type refreshTokenService struct {
	db               *gorm.DB
	refreshTokenRepo repository.RefreshTokenRepository
}

type RefreshTokenServiceConfig struct {
	DB               *gorm.DB
	RefreshTokenRepo repository.RefreshTokenRepository
}

func NewRefreshTokenService(c *RefreshTokenServiceConfig) RefreshTokenService {
	return &refreshTokenService{
		db:               c.DB,
		refreshTokenRepo: c.RefreshTokenRepo,
	}
}

func (r *refreshTokenService) CheckIfTokenExist(token string) (bool, uint, error) {
	tx := r.db.Begin()
	result, userID, err := r.refreshTokenRepo.CheckIfTokenExist(tx, token)
	if err != nil {
		tx.Rollback()
		return false, 0, err
	}

	tx.Commit()
	return result, userID, nil
}
