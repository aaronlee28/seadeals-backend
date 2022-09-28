package model

import "time"

type Transaction struct {
	Id            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	VoucherID     *uint     `json:"voucher_id"`
	Total         float64   `json:"total"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
