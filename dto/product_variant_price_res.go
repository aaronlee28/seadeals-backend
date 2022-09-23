package dto

type ProductVariantPriceRes struct {
	ProductID        uint                          `json:"product_id" binding:"required"`
	ProductName      string                        `json:"product_name" binding:"required"`
	ProductPromotion uint                          `json:"product_promotion" binding:"required"`
	ProductVariant   []*ProductVariantPromotionRes `json:"product_variant" binding:"required"`
}
