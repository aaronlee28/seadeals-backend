package dto

import "time"

type TransactionDetailsRes struct {
	VoucherID     uint      `json:"voucher_id"`
	Total         float64   `json:"Total"`
	PaymentType   string    `json:"payment_type"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
