package dto

import (
	"seadeals-backend/model"
	"time"
)

type TransactionDetailsRes struct {
	Id            uint      `json:"id"`
	VoucherID     uint      `json:"voucher_id"`
	Total         float64   `json:"total"`
	PaymentType   string    `json:"payment_type"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (_ *TransactionDetailsRes) FromTransaction(t *model.Transaction) *TransactionDetailsRes {
	return &TransactionDetailsRes{
		Id:            t.Id,
		VoucherID:     t.VoucherID,
		Total:         t.Total,
		PaymentType:   t.PaymentType,
		PaymentMethod: t.PaymentMethod,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}
