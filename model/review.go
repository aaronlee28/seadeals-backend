package model

import "gorm.io/gorm"

type ReviewQueryParam struct {
	Sort   string
	SortBy string
	Limit  int
	Page   int
}

type Review struct {
	gorm.Model  `json:"-"`
	ID          uint     `json:"id" gorm:"primaryKey"`
	UserID      uint     `json:"user_id"`
	User        *User    `json:"user"`
	ProductID   uint     `json:"product_id"`
	Product     *Product `json:"product"`
	Rating      int      `json:"rating"`
	ImageURL    *string  `json:"image_url"`
	Description *string  `json:"description"`
}
