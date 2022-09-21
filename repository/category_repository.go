package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
)

type ProductCategoryRepository interface {
	FindAllProductCategories(tx *gorm.DB) ([]*model.ProductCategory, error)
}

type productCategoryRepository struct{}

func NewProductCategoryRepository() ProductCategoryRepository {
	return &productCategoryRepository{}
}

func (r *productCategoryRepository) FindAllProductCategories(tx *gorm.DB) ([]*model.ProductCategory, error) {
	var categories []*model.ProductCategory

	result := tx.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}
