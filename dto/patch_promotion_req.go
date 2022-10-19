package dto

import "time"

type PatchPromotionReq struct {
	PromotionID uint      `json:"promotion_id"`
	Name        string    `json:"name"`
	Description string    `json:"Description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Quota       uint      `json:"quota"`
	MaxOrder    uint      `json:"max_order"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
}
