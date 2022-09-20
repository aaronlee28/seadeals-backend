package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	IconURL    string `json:"icon_url"`
	ParentID   uint   `json:"parent_id"`
}
