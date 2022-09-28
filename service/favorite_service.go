package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type FavoriteService interface {
	FavoriteToProduct(userID uint, productID uint) (*model.Favorite, error)
}

type favoriteService struct {
	db                 *gorm.DB
	favoriteRepository repository.FavoriteRepository
}

type FavoriteServiceConfig struct {
	DB                 *gorm.DB
	FavoriteRepository repository.FavoriteRepository
}

func NewFavoriteService(c *FavoriteServiceConfig) FavoriteService {
	return &favoriteService{
		db:                 c.DB,
		favoriteRepository: c.FavoriteRepository,
	}
}

func (f *favoriteService) FavoriteToProduct(userID uint, productID uint) (*model.Favorite, error) {
	tx := f.db.Begin()

	favorite, err := f.favoriteRepository.FavoriteToProduct(tx, userID, productID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return favorite, nil
}
