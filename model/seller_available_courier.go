package model

import "gorm.io/gorm"

type SellerAvailableCourier struct {
	gorm.Model `json:"-"`
	ID         uint     `json:"id" gorm:"primaryKey"`
	SellerID   uint     `json:"seller_id"`
	Seller     *Seller  `json:"seller"`
	CourierID  uint     `json:"courier_id"`
	Courier    *Courier `json:"courier"`
	IsSelected *bool    `json:"is_selected"`
	SlaDay     int      `json:"sla_day"`
}
