package dto

import "time"

type CheckoutCartRes struct {
	UserID        uint      `json:"user_id" binding:"required"`
	TransactionID uint      `json:"transaction_id" binding:"required"`
	Total         float64   `json:"total" binding:"required"`
	PaymentMethod string    `json:"payment_method" binding:"required"`
	CreatedAt     time.Time `json:"created_at" binding:"required"`
}
