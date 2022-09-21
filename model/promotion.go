package model

import (
	"gorm.io/gorm"
	"time"
)

type Promotion struct {
	gorm.Model  `json:"-"`
	ID          uint      `json:"id" gorm:"primaryKey"`
	ProductID   uint      `json:"product_id"`
	Product     *Product  `json:"product"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Quota       int       `json:"quota"`
	MaxOrder    int       `json:"max_order"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
}
