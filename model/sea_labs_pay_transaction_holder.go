package model

import "gorm.io/gorm"

type SeaLabsPayTransactionHolder struct {
	gorm.Model        `json:"-"`
	ID                uint         `json:"id" gorm:"primaryKey"`
	UserID            uint         `json:"user_id"`
	User              *User        `json:"user"`
	TransactionID     uint         `json:"transaction_id"`
	Transaction       *Transaction `json:"transaction"`
	TxnID             uint         `json:"txn_id"`
	Total             float64      `json:"total"`
	Sign              string       `json:"sign"`
	TransactionStatus *string      `json:"transaction_status"`
}
