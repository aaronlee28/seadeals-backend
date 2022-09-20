package service

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductService interface {
	FindProductDetailBySlug(slug string) (*model.Product, error)
}

type productService struct {
	db          *gorm.DB
	productRepo repository.ProductRepository
}

type ProductConfig struct {
	DB          *gorm.DB
	ProductRepo repository.ProductRepository
}

func NewProductService(config *ProductConfig) ProductService {
	return &productService{db: config.DB, productRepo: config.ProductRepo}
}

func (s *productService) FindProductDetailBySlug(slug string) (*model.Product, error) {
	tx := s.db.Begin()

	product, err := s.productRepo.FindProductBySlug(tx, slug)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	product, err = s.productRepo.FindProductDetailByID(tx, product.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return product, nil
}
