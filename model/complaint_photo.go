package model

import "gorm.io/gorm"

type ComplaintPhoto struct {
	gorm.Model  `json:"-"`
	ID          uint       `json:"id" gorm:"primaryKey"`
	ComplaintID uint       `json:"complaint_id"`
	Complaint   *Complaint `json:"complaint"`
	PhotoURL    string     `json:"photo_url"`
	PhotoName   string     `json:"photo_name"`
}
