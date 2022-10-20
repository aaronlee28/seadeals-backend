package model

import "gorm.io/gorm"

type Delivery struct {
	gorm.Model       `json:"-"`
	ID               uint                `json:"id" gorm:"primaryKey"`
	Address          string              `json:"address"`
	Status           string              `json:"status"`
	DeliveryNumber   string              `json:"delivery_number"`
	Total            float64             `json:"total"`
	Eta              int                 `json:"eta"`
	Weight           uint                `json:"weight"`
	CityDestination  string              `json:"city_destination"`
	OrderID          uint                `json:"order_id"`
	Order            *Order              `json:"order"`
	CourierID        uint                `json:"courier_id"`
	Courier          *Courier            `json:"courier"`
	DeliveryActivity []*DeliveryActivity `json:"delivery_activity"`
}

func (a Delivery) TableName() string {
	return "deliveries"
}
