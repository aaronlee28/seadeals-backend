package model

import "gorm.io/gorm"

type UserAddress struct {
	gorm.Model `json:"-"`
	ID         uint     `json:"id" gorm:"primaryKey"`
	AddressID  uint     `json:"address_id"`
	Address    *Address `json:"address"`
	UserID     uint     `json:"user_id"`
	User       *User    `json:"user"`
	IsMain     bool     `json:"is_main"`
}

func (u UserAddress) TableName() string {
	return "user_addresses"
}
