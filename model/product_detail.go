package model

import "gorm.io/gorm"

type ProductDetail struct {
	gorm.Model      `json:"-"`
	ID              uint    `json:"id" gorm:"primaryKey"`
	ProductID       uint    `json:"product_id"`
	Description     string  `json:"description"`
	VideoURL        string  `json:"video_url"`
	IsHazardous     bool    `json:"is_hazardous"`
	ConditionStatus string  `json:"condition_status"`
	Length          float64 `json:"length"`
	Width           float64 `json:"width"`
	Height          float64 `json:"height"`
	Weight          float64 `json:"weight"`
}
