package dto

import (
	"seadeals-backend/model"
	"time"
)

type SearchedProductRes struct {
	ProductID   uint      `json:"product_id" binding:"required"`
	ProductName string    `json:"product_name"`
	Slug        string    `json:"slug" binding:"required"`
	MediaURL    string    `json:"media_url" binding:"required"`
	MinPrice    uint      `json:"min_price" binding:"required"`
	MaxPrice    uint      `json:"max_price" binding:"required"`
	PromoPrice  float64   `json:"promo_price" binding:"required"`
	Rating      float64   `json:"rating" binding:"required"`
	Views       int       `json:"views" binding:"required"`
	Bought      int       `json:"bought" binding:"required"`
	City        string    `json:"city" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	UpdatedAt   time.Time `json:"updated_at" binding:"required"`
}

func (_ *SearchedProductRes) FromProduct(t *model.Product) *SearchedProductRes {
	return &SearchedProductRes{
		ProductID:   t.ID,
		ProductName: t.Name,
		Slug:        t.Slug,
		MediaURL:    "",
		MinPrice:    0,
		MaxPrice:    0,
		PromoPrice:  0,
		Rating:      0,
		Bought:      0,
		City:        "",
		Category:    "",
		UpdatedAt:   t.UpdatedAt,
	}
}
