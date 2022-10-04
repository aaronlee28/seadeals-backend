package dto

import "time"

type PostVoucherReq struct {
	Name        string    `json:"name" binding:"required"`
	Code        string    `json:"code" binding:"required,alphanum,min=1,max=5"`
	StartDate   time.Time `json:"start_date" binding:"required,ltefield=EndDate,gte"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Quota       int       `json:"quota" binding:"required,gt=0"`
	AmountType  string    `json:"amount_type" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,numeric,gt=0"`
	MinSpending float64   `json:"min_spending" binding:"numeric,gte=0"`
}
