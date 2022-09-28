package dto

import "time"

type PatchVoucherReq struct {
	Name        string    `json:"name"`
	StartDate   time.Time `json:"start_date" binding:"omitempty,ltefield=EndDate,gte"`
	EndDate     time.Time `json:"end_date"`
	Quota       int       `json:"quota" binding:"omitempty,gt=0"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount" binding:"omitempty,numeric,gt=0"`
	MinSpending float64   `json:"min_spending" binding:"omitempty,gte=0"`
}
