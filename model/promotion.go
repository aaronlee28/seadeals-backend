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
	SellerID    uint      `json:"seller_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Quota       uint      `json:"quota"`
	MaxOrder    uint      `json:"max_order"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
	BannerURL   string    `json:"banner_url"`
}
