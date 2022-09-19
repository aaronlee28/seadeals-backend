package model

import "gorm.io/gorm"

type SubDistrict struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	DistrictID string    `json:"district_id"`
	District   *District `json:"district"`
	Name       string    `json:"name"`
}
