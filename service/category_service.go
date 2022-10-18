package service

import (
	"gorm.io/gorm"
	"seadeals-backend/helper"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductCategoryService interface {
	FindCategories(query *model.CategoryQuery) ([]*model.ProductCategory, int64, int64, error)
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

func (s *productCategoryService) FindCategories(query *model.CategoryQuery) ([]*model.ProductCategory, int64, int64, error) {
	tx := s.db.Begin()
	var err error
	defer helper.CommitOrRollback(tx, &err)

	categories, totalPage, totalData, err := s.productCategoryRepository.FindCategories(tx, query)
	if err != nil {
		return nil, 0, 0, err
	}

	return categories, totalPage, totalData, nil
}
