package service

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
	"seadeals-backend/repository"
)

type ProductCategoryService interface {
	FindAllProductCategories(query *model.CategoryQuery) ([]*model.ProductCategory, int64, int64, error)
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

func (s *productCategoryService) FindAllProductCategories(query *model.CategoryQuery) ([]*model.ProductCategory, int64, int64, error) {
	tx := s.db.Begin()

	categories, totalPage, totalData, err := s.productCategoryRepository.FindAllProductCategories(tx, query)
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err.Error())
	}

	tx.Commit()
	return categories, totalPage, totalData, nil
}
