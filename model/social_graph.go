package model

import "gorm.io/gorm"

type SocialGraph struct {
	gorm.Model `json:"-"`
	ID         uint    `json:"id" gorm:"primaryKey"`
	IsFollow   bool    `json:"is_follow"`
	UserID     uint    `json:"user_id"`
	User       *User   `json:"user"`
	SellerID   uint    `json:"seller_id"`
	Seller     *Seller `json:"seller"`
}
