package dto

import "time"

type CreateGlobalVoucher struct {
	Name        string    `json:"name" binding:"required"`
	Code        string    `json:"code" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Quota       uint      `json:"quota" binding:"required"`
	AmountType  string    `json:"amount_type" binding:"required"`
	Amount      float64   `json:"amount" binding:"required"`
	MinSpending float64   `json:"min_spending" binding:"required"`
}
