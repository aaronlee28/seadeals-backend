package model

import (
	"gorm.io/gorm"
	"strconv"
)

type Seller struct {
	gorm.Model  `json:"-"`
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	UserID      uint         `json:"user_id"`
	User        *User        `json:"user"`
	Description string       `json:"description"`
	AddressID   uint         `json:"address_id"`
	Address     *Address     `json:"address"`
	PictureURL  string       `json:"picture_url"`
	BannerURL   string       `json:"banner_url"`
	SocialGraph *SocialGraph `json:"following"`
	AllowPrint  bool         `json:"allow_print"`
}

func (u *Seller) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(u).Update("slug", u.Name+"."+strconv.FormatUint(uint64(u.ID), 10))
	return
}
