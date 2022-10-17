package model

import "gorm.io/gorm"

const (
	SortReviewDefault   = "desc"
	SortByReviewDefault = ""
	LimitReviewDefault  = 99999999
	PageReviewDefault   = 1
)

type ReviewQueryParam struct {
	Sort                string
	SortBy              string
	Limit               uint
	Page                uint
	Rating              uint
	WithImageOnly       bool
	WithDescriptionOnly bool
}

type Review struct {
	gorm.Model  `json:"-"`
	ID          uint     `json:"id" gorm:"primaryKey"`
	UserID      uint     `json:"user_id"`
	User        *User    `json:"user"`
	ProductID   uint     `json:"product_id"`
	Product     *Product `json:"product"`
	Rating      int      `json:"rating"`
	ImageURL    *string  `json:"image_url"`
	ImageName   *string  `json:"image_name"`
	Description *string  `json:"description"`
}
