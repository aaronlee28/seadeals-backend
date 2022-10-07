package dto

import "time"

type CreateGlobalVoucher struct {
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Quota       uint      `json:"quota"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
	MinSpending float64   `json:"min_spending"`
}
