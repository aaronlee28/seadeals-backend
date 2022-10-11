package model

import "gorm.io/gorm"

type Notification struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id"`
	SellerID   uint   `json:"seller_id"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
}
