package dto

import "seadeals-backend/model"

type PostCreateProductRes struct {
	Product              *model.Product                `json:"product"`
	ProductDetail        *model.ProductDetail          `json:"product_detail"`
	ProductPhoto         []*model.ProductPhoto         `json:"product_photo"`
	ProductVariantDetail []*model.ProductVariantDetail `json:"product_variant_detail_1"`
}
