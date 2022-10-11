package dto

import (
	"seadeals-backend/model"
	"time"
)

type TransactionsRes struct {
	Id            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	VoucherID     *uint     `json:"voucher_id"`
	Total         float64   `json:"total"`
	PaymentType   string    `json:"payment_type"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (_ *TransactionsRes) FromTransaction(t *model.Transaction) *TransactionsRes {
	return &TransactionsRes{
		Id:            t.ID,
		UserID:        t.UserID,
		VoucherID:     t.VoucherID,
		Total:         t.Total,
		PaymentMethod: t.PaymentMethod,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}

const (
	TransactionWaitingPayment = "waiting for payment"
	TransactionPayed          = "payed"
)
