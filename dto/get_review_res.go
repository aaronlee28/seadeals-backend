package dto

import "seadeals-backend/model"

type GetReviewRes struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	UserID        uint    `json:"user_id"`
	ProductID     uint    `json:"product_id"`
	UserUsername  string  `json:"username"`
	UserAvatarURL *string `json:"avatar_url"`
	Rating        int     `json:"rating"`
	Description   string  `json:"description"`
}

func (_ *GetReviewRes) From(r *model.Review) *GetReviewRes {
	return &GetReviewRes{
		ID:            r.ID,
		UserID:        r.UserID,
		ProductID:     r.ProductID,
		UserUsername:  r.User.Username,
		UserAvatarURL: r.User.AvatarURL,
		Rating:        r.Rating,
		Description:   r.Description,
	}
}
