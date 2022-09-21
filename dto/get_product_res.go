package dto

import "seadeals-backend/model"

type ProductRes struct {
	MinPrice float64        `json:"min_price"`
	MaxPrice float64        `json:"max_price"`
	Product  *GetProductRes `json:"product"`
}

type GetProductRes struct {
	ID            uint    `json:"id"`
	Price         float64 `json:"price"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	PictureURL    string  `json:"picture_url"`
	City          string  `json:"city"`
	Rating        float64 `json:"rating"`
	TotalReviewer int64   `json:"total_reviewer"`
	TotalSold     uint    `json:"totalSold"`
}

type SellerProductsCustomTable struct {
	Min                        float64 `json:"min"`
	Max                        float64 `json:"max"`
	ProductID                  uint    `json:"product_id"`
	model.ProductVariantDetail `json:"product_variant_detail"`
}

func (_ SellerProductsCustomTable) TableName() string {
	return "product_variant_details"
}
