package dto

import "seadeals-backend/model"

type FavoriteRes struct {
	ID            uint          `json:"id"`
	UserID        uint          `json:"user_id"`
	ProductID     uint          `json:"product_id"`
	SellerID      uint          `json:"seller_id"`
	Seller        *model.Seller `json:"seller"`
	Name          string        `json:"name"`
	Slug          string        `json:"slug"`
	ProductPhotos *ProductPhoto `json:"photo_url"`
}
