package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model `json:"-"`
	ID         uint     `json:"id" gorm:"primaryKey"`
	IsFavorite bool     `json:"is_favorite"`
	UserID     uint     `json:"user_id"`
	User       *User    `json:"user"`
	ProductID  uint     `json:"product_id"`
	Product    *Product `json:"product"`
}
