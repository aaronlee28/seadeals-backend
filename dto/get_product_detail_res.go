package dto

import "seadeals-backend/model"

type ProductDetailRes struct {
	TotalStock    uint `json:"total_stock"`
	model.Product `json:"product"`
}

func (_ ProductDetailRes) TableName() string {
	return "products"
}
