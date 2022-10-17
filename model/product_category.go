package model

import "gorm.io/gorm"

type ProductCategory struct {
	gorm.Model `json:"-"`
	ID         uint               `json:"id" gorm:"primaryKey"`
	Name       string             `json:"name"`
	Slug       string             `json:"slug"`
	IconURL    string             `json:"icon_url"`
	ParentID   *uint              `json:"parent_id" gorm:"foreignKey:ID"`
	Children   []*ProductCategory `json:"children" gorm:"foreignKey:ParentID"`
}

type CategoryQuery struct {
	Search   string
	Limit    string
	Page     string
	SellerID uint
	ParentID uint
	FindAll  bool
}

func (a ProductCategory) TableName() string {
	return "product_categories"
}
