package model

import (
	"gorm.io/gorm"
	"time"
)

const PercentageType = "percentage"
const NominalType = "nominal"

type Voucher struct {
	gorm.Model  `json:"-"`
	ID          uint      `json:"id" gorm:"primaryKey"`
	SellerID    uint      `json:"seller_id"`
	Seller      *Seller   `json:"seller"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Quota       int       `json:"quota"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
	MinSpending float64   `json:"min_spending"`
}
