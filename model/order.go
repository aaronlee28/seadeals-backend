package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model    `json:"-"`
	ID            uint         `json:"id" gorm:"primaryKey"`
	SellerID      uint         `json:"seller_id"`
	Seller        *Seller      `json:"seller"`
	VoucherID     *uint        `json:"voucher_id"`
	Voucher       *Voucher     `json:"voucher"`
	TransactionID uint         `json:"transaction_id"`
	Transaction   *Transaction `json:"transaction"`
	UserID        uint         `json:"user_id"`
	User          *User        `json:"user"`
	Total         float64      `json:"total"`
	Status        string       `json:"status"`
	OrderItems    []*OrderItem `json:"order_items"`
	Delivery      *Delivery    `json:"delivery"`
	Complaint     *Complaint   `json:"complaint"`
}
