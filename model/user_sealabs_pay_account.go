package model

import (
	"gorm.io/gorm"
	"time"
)

type UserSealabsPayAccount struct {
	gorm.Model    `json:"-"`
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name"`
	UserID        uint      `json:"user_id"`
	User          *User     `json:"user"`
	ActiveDate    time.Time `json:"active_date"`
	AccountNumber string    `json:"account_number"`
	IsMain        bool      `json:"is_main"`
}
