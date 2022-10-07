package dto

import "seadeals-backend/model"

type ProductDetailRes struct {
	TotalStock    uint    `json:"total_stock"`
	AverageRating float64 `json:"average_rating"`
	RatingCount   uint    `json:"rating_count"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	model.Product `json:"product"`
}

func (_ ProductDetailRes) TableName() string {
	return "products"
}
