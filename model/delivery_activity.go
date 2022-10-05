package model

import "gorm.io/gorm"

type DeliveryActivity struct {
	gorm.Model  `json:"-"`
	ID          uint      `json:"id" gorm:"primaryKey"`
	Description string    `json:"description"`
	DeliveryID  uint      `json:"delivery_id"`
	Delivery    *Delivery `json:"delivery"`
}

func (a DeliveryActivity) TableName() string {
	return "delivery_activities"
}
