package dto

import "seadeals-backend/model"

type ProductRes struct {
	MinPriceBeforeDisc float64        `json:"min_price_before_disc"`
	MaxPriceBeforeDisc float64        `json:"max_price_before_disc"`
	MinPrice           float64        `json:"min_price"`
	MaxPrice           float64        `json:"max_price"`
	Product            *GetProductRes `json:"product"`
}

type GetProductRes struct {
	ID              uint    `json:"id"`
	Price           float64 `json:"price"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	MediaURL        string  `json:"media_url"`
	City            string  `json:"city"`
	Rating          float64 `json:"rating"`
	TotalReviewer   int64   `json:"total_reviewer"`
	TotalSold       uint    `json:"total_sold"`
	PromotionAmount float64 `json:"promotion_amount"`
}

type SellerProductsCustomTable struct {
	MinBeforeDisc   float64 `json:"min_before_disc"`
	MaxBeforeDisc   float64 `json:"max_before_disc"`
	Min             float64 `json:"min"`
	Max             float64 `json:"max"`
	Avg             float64 `json:"review_avg"`
	Count           int64   `json:"review_count"`
	ProductID       uint    `json:"product_id"`
	PromotionAmount float64 `json:"promotion_amount"`
	model.Product   `json:"product"`
}

func (_ SellerProductsCustomTable) TableName() string {
	return "products"
}
