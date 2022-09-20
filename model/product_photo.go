package model

import "gorm.io/gorm"

type ProductPhoto struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	ProductID  uint   `json:"product_id"`
	PhotoURL   string `json:"photo_url"`
	Name       string `json:"name"`
}
