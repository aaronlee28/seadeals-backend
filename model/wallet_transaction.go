package model

import (
	"gorm.io/gorm"
)

type WalletTransaction struct {
	gorm.Model    `json:"-"`
	ID            uint    `json:"id" gorm:"primaryKey"`
	WalletID      uint    `json:"wallet_id"`
	Wallet        *Wallet `json:"wallet"`
	Total         float64 `json:"total"`
	PaymentMethod string  `json:"payment_method"`
	PaymentType   string  `json:"payment_type"`
	Description   string  `json:"description"`
}
