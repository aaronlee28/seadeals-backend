package model

import "gorm.io/gorm"

type Review struct {
	gorm.Model  `json:"-"`
	ID          uint     `json:"id" gorm:"primaryKey"`
	UserID      uint     `json:"user_id"`
	User        *User    `json:"user"`
	ProductID   uint     `json:"product_id"`
	Product     *Product `json:"product"`
	Rating      int      `json:"rating"`
	Description string   `json:"description"`
}
