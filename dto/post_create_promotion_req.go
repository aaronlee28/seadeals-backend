package dto

import "time"

type CreatePromotionReq struct {
	ProductID   uint      `json:"product_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required,ltefield=EndDate,gte"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	Quota       uint      `json:"quota" binding:"required,numeric,gte=1"`
	MaxOrder    uint      `json:"max_order" binding:"required,numeric,gte=1"`
	AmountType  string    `json:"amount_type" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,numeric,gte=1"`
	BannerURL   string    `json:"banner_url" binding:"required"`
}
