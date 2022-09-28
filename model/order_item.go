package model

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model             `json:"-"`
	ID                     uint                  `json:"id" gorm:"primaryKey"`
	ProductVariantDetailID uint                  `json:"product_variant_detail_id"`
	ProductVariantDetail   *ProductVariantDetail `json:"product_variant_detail"`
	PromotionID            *uint                 `json:"promotion_id"`
	Promotion              *Promotion            `json:"promotion"`
	OrderID                *uint                 `json:"order_id"`
	Order                  *Order                `json:"order"`
	UserID                 uint                  `json:"user_id"`
	User                   *User                 `json:"user"`
	Quantity               uint                  `json:"quantity"`
	Subtotal               float64               `json:"subtotal"`
}
