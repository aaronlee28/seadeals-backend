package model

import "gorm.io/gorm"

type SeaLabsPayTopUpHolder struct {
	gorm.Model        `json:"-"`
	ID                uint    `json:"id" gorm:"primaryKey"`
	UserID            uint    `json:"user_id"`
	User              *User   `json:"user"`
	TxnID             uint    `json:"txn_id"`
	Total             float64 `json:"total"`
	Sign              string  `json:"sign"`
	TransactionStatus *string `json:"transaction_status"`
}
