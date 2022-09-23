package model

import "gorm.io/gorm"

type ReviewPhoto struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	ReviewID   uint   `json:"review_id"`
	PhotoURL   string `json:"photo_url"`
	Name       string `json:"name"`
}
