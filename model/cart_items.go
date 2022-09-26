package model

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model             `json:"-"`
	ID                     uint                  `json:"id" gorm:"primaryKey"`
	ProductVariantDetailID uint                  `json:"product_variant_detail_id"`
	ProductVariantDetail   *ProductVariantDetail `json:"product_variant_detail"`
	UserID                 uint                  `json:"user_id"`
	User                   *User                 `json:"user"`
	Quantity               uint                  `json:"quantity"`
}
