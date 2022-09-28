package dto

import "time"

type GetPromotionRes struct {
	ID              uint      `json:"id"`
	ProductID       uint      `json:"product_id"`
	Name            string    `json:"name"`
	Description     string    `json:"Description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	AmountType      string    `json:"amount_type"`
	Amount          float64   `json:"amount"`
	ProductPhotoURL string    `json:"product_photo_url"`
	Status          string    `json:"status"`
}
