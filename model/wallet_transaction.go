package model

import (
	"gorm.io/gorm"
	"time"
)

type WalletTransaction struct {
	gorm.Model    `json:"-"`
	ID            uint         `json:"id" gorm:"primaryKey"`
	WalletID      uint         `json:"wallet_id"`
	Wallet        *Wallet      `json:"wallet"`
	TransactionID *uint        `json:"transaction_id"`
	Transaction   *Transaction `json:"transaction"`
	Total         float64      `json:"amount"`
	PaymentMethod string       `json:"payment_method"`
	PaymentType   string       `json:"payment_type"`
	Description   string       `json:"description"`
	CreatedAt     time.Time    `json:"created_at"`
}
