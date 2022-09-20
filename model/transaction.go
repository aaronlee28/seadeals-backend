package model

import "time"

type transaction struct {
	Id        uint      `json:"id"`
	StatusID  uint      `json:"status_id"`
	OrderID   uint      `json:"order_id"`
	CourierID uint      `json:"courier_id"`
	AddressID uint      `json:"address_id"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
