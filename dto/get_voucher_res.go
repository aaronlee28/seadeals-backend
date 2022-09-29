package dto

import (
	"seadeals-backend/model"
	"time"
)

const (
	statusUpcoming = "upcoming"
	statusOnGoing  = "ongoing"
	statusEnded    = "ended"
)

type GetVouchersRes struct {
	Limit         uint             `json:"limit"`
	Page          uint             `json:"page"`
	TotalPages    uint             `json:"total_pages"`
	TotalVouchers uint             `json:"total_vouchers"`
	Vouchers      []*GetVoucherRes `json:"vouchers"`
}

type GetVoucherRes struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SellerID    uint      `json:"seller_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status"`
	Quota       int       `json:"quota"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
	MinSpending float64   `json:"min_spending"`
}

func (_ *GetVoucherRes) From(v *model.Voucher) *GetVoucherRes {
	status := statusOnGoing
	if time.Now().After(v.EndDate) {
		status = statusEnded
	} else if v.StartDate.After(time.Now()) {
		status = statusUpcoming
	}

	return &GetVoucherRes{
		ID:          v.ID,
		SellerID:    v.SellerID,
		Name:        v.Name,
		Code:        v.Code,
		StartDate:   v.StartDate,
		EndDate:     v.EndDate,
		Status:      status,
		Quota:       v.Quota,
		AmountType:  v.AmountType,
		Amount:      v.Amount,
		MinSpending: v.MinSpending,
	}
}
