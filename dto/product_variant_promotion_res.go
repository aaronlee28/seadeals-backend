package dto

type ProductVariantPromotionRes struct {
	VariantID           uint   `json:"variant_id" binding:"required"`
	VariantName         string `json:"variant_name" binding:"required"`
	Price               uint   `json:"price" binding:"required"`
	PriceAfterPromotion uint   `json:"price_after_promotion" binding:"required"`
}
