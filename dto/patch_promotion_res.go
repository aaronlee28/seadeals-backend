package dto

import (
	"seadeals-backend/model"
	"time"
)

type PatchPromotionRes struct {
	Name        string    `json:"name"`
	Description string    `json:"Description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Quota       uint      `json:"quota"`
	MaxOrder    uint      `json:"max_order"`
	AmountType  string    `json:"amount_type"`
	Amount      float64   `json:"amount"`
}

func (_ *PatchPromotionRes) PatchFromPromotion(t *model.Promotion) *PatchPromotionRes {
	return &PatchPromotionRes{
		Name:        t.Name,
		Description: t.Description,
		StartDate:   t.StartDate,
		EndDate:     t.EndDate,
		Quota:       t.Quota,
		MaxOrder:    t.MaxOrder,
		AmountType:  t.AmountType,
		Amount:      t.Amount,
	}
}
