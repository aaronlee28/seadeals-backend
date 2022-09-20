package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model    `json:"-"`
	ID            uint         `json:"id" gorm:"primaryKey"`
	Address       string       `json:"address"`
	Zipcode       string       `json:"zipcode"`
	SubDistrictID uint         `json:"sub_district_id"`
	SubDistrict   *SubDistrict `json:"sub_district"`
	UserID        uint         `json:"user_id"`
	User          *User        `json:"user"`
}

func (a Address) TableName() string {
	return "addresses"
}
