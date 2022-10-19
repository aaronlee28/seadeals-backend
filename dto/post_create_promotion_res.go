package dto

import "time"

type CreatePromotionRes struct {
	ID          uint      `json:"id"`
	ProductID   uint      `json:"product_id"`
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
