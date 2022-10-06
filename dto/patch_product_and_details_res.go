package dto

import "seadeals-backend/model"

type PatchProductAndDetailsRes struct {
	Product       *model.Product       `json:"product" binding:"omitempty"`
	ProductDetail *model.ProductDetail `json:"product_detail" binding:"omitempty"`
}
