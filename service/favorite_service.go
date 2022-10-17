package service

import (
	"gorm.io/gorm"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type FavoriteService interface {
	FavoriteToProduct(userID uint, productID uint) (*model.Favorite, uint, error)
	GetUserFavoriteCount(userID uint) (uint, error)
}

type favoriteService struct {
	db                 *gorm.DB
	favoriteRepository repository.FavoriteRepository
	productRepository  repository.ProductRepository
}

type FavoriteServiceConfig struct {
	DB                 *gorm.DB
	FavoriteRepository repository.FavoriteRepository
	ProductRepository  repository.ProductRepository
}

func NewFavoriteService(c *FavoriteServiceConfig) FavoriteService {
	return &favoriteService{
		db:                 c.DB,
		favoriteRepository: c.FavoriteRepository,
		productRepository:  c.ProductRepository,
	}
}

func (f *favoriteService) FavoriteToProduct(userID uint, productID uint) (*model.Favorite, uint, error) {
	tx := f.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	favorite, err := f.favoriteRepository.FavoriteToProduct(tx, userID, productID)
	if err != nil {
		return nil, 0, err
	}

	product, err := f.productRepository.UpdateProductFavoriteCount(tx, productID, favorite.IsFavorite)
	if err != nil {
		return nil, 0, err
	}

	return favorite, product.FavoriteCount, nil
}

func (f *favoriteService) GetUserFavoriteCount(userID uint) (uint, error) {
	tx := f.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	userFavCount, err := f.favoriteRepository.GetUserFavoriteCount(tx, userID)
	if err != nil {
		return 0, err
	}

	return userFavCount, nil
}
