package model

import (
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	gorm.Model   `json:"-"`
	ID           uint       `json:"id" gorm:"primaryKey"`
	UserID       uint       `json:"user_id"`
	User         User       `json:"-"`
	Balance      float64    `json:"balance"`
	Pin          *string    `json:"pin"`
	Status       string     `json:"status"`
	BlockedUntil *time.Time `json:"blocked_until"`
}

const WalletActive = "active"
