package dto

import "seadeals-backend/model"

type CreateComplaintReq struct {
	OrderID uint `json:"order_id" binding:"required"`
	Photos  []struct {
		PhotoURL  string `json:"photo_url" binding:"required"`
		PhotoName string `json:"photo_name" binding:"required"`
	} `json:"photos"`
	Description string `json:"description" binding:"required"`
}

type CreateComplaintRes struct {
	Order           *model.Order            `json:"order"`
	ComplaintPhotos []*model.ComplaintPhoto `json:"complaint_photos"`
	Description     string                  `json:"description"`
}
