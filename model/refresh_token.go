package model

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id"`
	User       User   `json:"user"`
	Token      string `json:"token"`
}
