package dto

import (
	"seadeals-backend/model"
	"time"
)

type TransactionDetailsRes struct {
	Id            uint           `json:"id"`
	VoucherID     *uint          `json:"voucher_id"`
	Voucher       *model.Voucher `json:"voucher"`
	Total         float64        `json:"total"`
	PaymentType   string         `json:"payment_type"`
	PaymentMethod string         `json:"payment_method"`
	Orders        []*model.Order `json:"orders"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
