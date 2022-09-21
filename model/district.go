package model

import "gorm.io/gorm"

type District struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	CityID     uint   `json:"city_id"`
	City       *City  `json:"city"`
	Name       string `json:"name"`
}
