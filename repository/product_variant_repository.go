package repository

import (
	"gorm.io/gorm"
	"seadeals-backend/apperror"
	"seadeals-backend/model"
)

type ProductVariantRepository interface {
	FindAllProductVariantByProductID(tx *gorm.DB, productID uint) ([]*model.ProductVariantDetail, error)
}

type productVariantRepository struct {
}

func NewProductVariantRepository() ProductVariantRepository {
	return &productVariantRepository{}
}

func (r *productVariantRepository) FindAllProductVariantByProductID(tx *gorm.DB, productID uint) ([]*model.ProductVariantDetail, error) {
	var productVariantDetails []*model.ProductVariantDetail
	result := tx.Where("product_id = ?", productID).Preload("ProductVariant1").Preload("ProductVariant2").Find(&productVariantDetails)
	if result.Error != nil {
		return nil, result.Error
	}
	if int(result.RowsAffected) == 0 {
		return nil, new(apperror.ProductNotFoundError)
	}
	return productVariantDetails, nil
}
