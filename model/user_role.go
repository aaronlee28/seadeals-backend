package model

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model `json:"-"`
	ID         uint `json:"id" gorm:"primaryKey"`
	UserID     uint `json:"user_id"`
	User       User `json:"user"`
	RoleID     uint `json:"role_id"`
	Role       Role `json:"role"`
}
