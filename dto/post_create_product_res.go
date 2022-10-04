package dto

import "seadeals-backend/model"

type PostCreateProductRes struct {
	Product              *model.Product              `json:"product"`
	ProductDetail        *model.ProductDetail        `json:"product_detail"`
	ProductPhoto         []*model.ProductPhoto       `json:"product_photo"`
	ProductVariant1      *model.ProductVariant       `json:"product_variant_1"`
	ProductVariant2      *model.ProductVariant       `json:"product_variant_2" binding:"omitempty"`
	ProductVariantDetail *model.ProductVariantDetail `json:"product_variant_detail_1"`
}
