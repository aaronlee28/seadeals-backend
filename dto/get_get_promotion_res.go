package dto

import (
	"seadeals-backend/model"
	"time"
)

type GetPromotionRes struct {
	ID              uint      `json:"id"`
	ProductID       uint      `json:"product_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	AmountType      string    `json:"amount_type"`
	Amount          float64   `json:"amount"`
	ProductPhotoURL string    `json:"product_photo_url"`
}

func (_ *GetPromotionRes) FromPromotion(t *model.Promotion) *GetPromotionRes {
	return &GetPromotionRes{
		ID:          t.ID,
		ProductID:   t.ProductID,
		Name:        t.Name,
		Description: t.Description,
		StartDate:   t.StartDate,
		EndDate:     t.EndDate,
		AmountType:  t.AmountType,
		Amount:      t.Amount,
	}
}
