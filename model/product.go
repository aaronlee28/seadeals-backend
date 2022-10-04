package model

import (
	"gorm.io/gorm"
	"strconv"
)

type Product struct {
	gorm.Model           `json:"-"`
	ID                   uint                    `json:"id" gorm:"primaryKey"`
	CategoryID           uint                    `json:"category_id"`
	Category             *ProductCategory        `json:"category"`
	SellerID             uint                    `json:"seller_id"`
	Seller               *Seller                 `json:"seller"`
	Name                 string                  `json:"name"`
	Slug                 string                  `json:"slug"`
	IsBulkEnabled        bool                    `json:"is_bulk_enabled"`
	SoldCount            int                     `json:"sold_count"`
	FavoriteCount        uint                    `json:"favorite_count"`
	IsArchived           bool                    `json:"is_archived"`
	ProductVariantDetail []*ProductVariantDetail `json:"product_variant_detail"`
	ProductDetail        *ProductDetail          `json:"product_detail"`
	ProductPhotos        []*ProductPhoto         `json:"product_photos"`
	Promotion            *Promotion              `json:"promotion"`
	Favorite             *Favorite               `json:"favorite"`
	MinQuantity          uint                    `json:"min_quantity"`
	MaxQuantity          uint                    `json:"max_quantity"`
}

type SellerProductQuery struct {
	SortBy string `json:"sortBy"`
	Sort   string `json:"sort"`
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

func (u *Product) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(u).Update("slug", u.Name+"."+strconv.FormatUint(uint64(u.ID), 10))
	return
}
