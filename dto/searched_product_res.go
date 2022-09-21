package dto

import "seadeals-backend/model"

type SearchedProductRes struct {
	ProductID  uint    `json:"product_id" binding:"required"`
	Slug       string  `json:"slug" binding:"required"`
	MediaURL   string  `json:"media_url" binding:"required"`
	MinPrice   uint    `json:"min_price" binding:"required"`
	MaxPrice   uint    `json:"max_price" binding:"required"`
	PromoPrice float64 `json:"promo_price" binding:"required"`
	Rating     uint    `json:"rating" binding:"required"`
	Bought     int     `json:"Bought" binding:"required"`
	City       string  `json:"city" binding:"required"`
}

func (_ *SearchedProductRes) FromProduct(t *model.Product) *SearchedProductRes {
	return &SearchedProductRes{
		ProductID:  t.ID,
		Slug:       t.Slug,
		MediaURL:   "",
		MinPrice:   0,
		MaxPrice:   0,
		PromoPrice: t.Promotion.Amount,
		Rating:     0,
		Bought:     t.SoldCount,
		City:       "",
	}
}
