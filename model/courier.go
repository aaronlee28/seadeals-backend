package model

import "gorm.io/gorm"

type Courier struct {
	gorm.Model  `json:"-"`
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
}
