package dto

import (
	"seadeals-backend/model"
	"time"
)

type GetReviewsRes struct {
	Limit         uint            `json:"limit"`
	Page          uint            `json:"page"`
	TotalPages    uint            `json:"total_pages"`
	TotalReviews  uint            `json:"total_reviews"`
	AverageRating float64         `json:"average_rating"`
	Reviews       []*GetReviewRes `json:"reviews"`
}

type GetReviewRes struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id"`
	ProductID     uint      `json:"product_id"`
	UserUsername  string    `json:"username"`
	UserAvatarURL *string   `json:"avatar_url"`
	Rating        int       `json:"rating"`
	ImageURL      *string   `json:"image_url"`
	ImageName     *string   `json:"image_name"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
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
		ImageURL:      r.ImageURL,
		ImageName:     r.ImageName,
		CreatedAt:     r.CreatedAt,
	}
}
