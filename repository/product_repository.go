package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/model"
)

type ProductRepository interface {
	FindProductDetailByID(tx *gorm.DB, id uint) (*model.Product, error)
	FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) FindProductDetailByID(tx *gorm.DB, id uint) (*model.Product, error) {
	var product *model.Product
	result := tx.Preload("ProductPhotos", "product_id = ?", id).Preload("ProductDetail", "product_id = ?", id).First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *productRepository) FindProductBySlug(tx *gorm.DB, slug string) (*model.Product, error) {
	var product *model.Product
	result := tx.First(&product, "slug = ?", slug)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
