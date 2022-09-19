package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Password   string    `json:"-"`
	AvatarURL  *string   `json:"avatar_url"`
	Gender     string    `json:"gender"`
	BirthDate  time.Time `json:"birth_date"`
}
