package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model    `json:"-"`
	ID            uint         `json:"id" gorm:"primaryKey"`
	Address       string       `json:"address"`
	Zipcode       string       `json:"zipcode"`
	SubDistrictID uint         `json:"sub_district_id"`
	SubDistrict   *SubDistrict `json:"sub_district"`
	UserAddress   *UserAddress `json:"user_address"`
}

func (a Address) TableName() string {
	return "addresses"
}
