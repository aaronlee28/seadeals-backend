package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
}
