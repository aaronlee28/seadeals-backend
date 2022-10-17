package model

import "gorm.io/gorm"

type AccountHolder struct {
	gorm.Model `json:"-"`
	ID         uint    `json:"id" gorm:"primaryKey"`
	UserID     uint    `json:"user_id"`
	User       *User   `json:"user"`
	OrderID    uint    `json:"order_id"`
	Order      *Order  `json:"order"`
	SellerID   uint    `json:"seller_id"`
	Seller     *Seller `json:"seller"`
	Total      float64 `json:"total"`
	HasTaken   bool    `json:"has_taken"`
}
