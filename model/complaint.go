package model

import "gorm.io/gorm"

type Complaint struct {
	gorm.Model      `json:"-"`
	ID              uint              `json:"id" gorm:"primaryKey"`
	OrderID         uint              `json:"order_id"`
	Order           *Order            `json:"order"`
	Description     string            `json:"description"`
	ComplaintPhotos []*ComplaintPhoto `json:"complaint_photos"`
}
