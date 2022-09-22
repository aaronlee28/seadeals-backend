package model

import "gorm.io/gorm"

type Seller struct {
	gorm.Model  `json:"-"`
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	UserID      uint     `json:"user_id"`
	User        *User    `json:"user"`
	Description string   `json:"description"`
	AddressID   uint     `json:"address_id"`
	Address     *Address `json:"address"`
	PictureURL  string   `json:"picture_url"`
	BannerURL   string   `json:"banner_url"`
}
