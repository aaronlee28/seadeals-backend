package dto

import "seadeals-backend/model"

type ProductVariantPromotionRes struct {
	VariantID           uint    `json:"variant_id" binding:"required"`
	Variant1Name        *string `json:"variant_1_name" binding:"required"`
	Variant2Name        *string `json:"variant_2_name" binding:"required"`
	Price               float64 `json:"price" binding:"required"`
	PriceAfterPromotion float64 `json:"price_after_promotion" binding:"required"`
}

func (_ *ProductVariantPromotionRes) FromProductVariantDetail(t model.ProductVariantDetail) *ProductVariantPromotionRes {
	return &ProductVariantPromotionRes{
		VariantID:           t.ID,
		Variant1Name:        t.Variant1Value,
		Variant2Name:        t.Variant2Value,
		Price:               t.Price,
		PriceAfterPromotion: 0,
	}
}
