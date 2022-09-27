package model

import "gorm.io/gorm"

type Address struct {
	gorm.Model  `json:"-"`
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id"`
	CityID      string `json:"city_id"`
	ProvinceID  string `json:"province_id"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Type        string `json:"type"`
	PostalCode  string `json:"postal_code"`
	SubDistrict string `json:"sub_district"`
	Address     string `json:"address"`
	IsMain      bool   `json:"is_main"`
}

func (a Address) TableName() string {
	return "addresses"
}
