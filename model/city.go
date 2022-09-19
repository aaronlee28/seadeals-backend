package model

import "gorm.io/gorm"

type City struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	ProvinceID string    `json:"province_id"`
	Province   *Province `json:"province"`
	Name       string    `json:"name"`
}

func (c City) TableName() string {
	return "cities"
}
