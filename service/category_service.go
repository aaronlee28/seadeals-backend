package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductCategoryService interface {
	FindAllProductCategories() ([]*model.ProductCategory, error)
}

type productCategoryService struct {
	db                        *gorm.DB
	productCategoryRepository repository.ProductCategoryRepository
}

type ProductCategoryServiceConfig struct {
	DB                        *gorm.DB
	ProductCategoryRepository repository.ProductCategoryRepository
}

func NewProductCategoryService(c *ProductCategoryServiceConfig) ProductCategoryService {
	return &productCategoryService{
		db:                        c.DB,
		productCategoryRepository: c.ProductCategoryRepository,
	}
}

func (s *productCategoryService) FindAllProductCategories() ([]*model.ProductCategory, error) {
	tx := s.db.Begin()

	categories, err := s.productCategoryRepository.FindAllProductCategories(tx)
	if err != nil {
		return nil, apperror.InternalServerError(err.Error())
	}

	tx.Commit()
	return categories, nil
}
