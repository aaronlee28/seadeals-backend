package model

import "time"

type Product struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	CategoryID    uint            `json:"category_id"`
	SellerID      uint            `json:"seller_id"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	IsBulkEnabled bool            `json:"is_bulk_enabled"`
	SoldCount     int             `json:"sold_count"`
	ViewsCount    uint            `json:"views_count"`
	IsArchived    bool            `json:"is_archived"`
	ProductDetail *ProductDetail  `json:"product_detail"`
	ProductPhotos []*ProductPhoto `json:"product_photos"`
	Promotion     *Promotion      `json:"promotion"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
