package dto

import (
	"seadeals-backend/model"
)

type SearchedProductRes struct {
	MinPrice       uint `json:"min_price" binding:"required"`
	MaxPrice       uint `json:"max_price" binding:"required"`
	*GetProductRes `json:"product"`
}

func (_ *SearchedProductRes) FromProduct(t *model.Product) *SearchedProductRes {
	return &SearchedProductRes{
		MinPrice: 0,
		MaxPrice: 0,
		GetProductRes: &GetProductRes{
			ID:            t.ID,
			Price:         0,
			Name:          t.Name,
			Slug:          t.Slug,
			MediaURL:      "",
			City:          "",
			Rating:        0,
			TotalReviewer: 0,
			TotalSold:     0,
		},
	}
}
