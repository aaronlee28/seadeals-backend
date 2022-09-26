package dto

type AddToCartReq struct {
	ProductVariantDetailID uint `json:"product_variant_detail_id" binding:"required"`
	UserID                 uint `json:"user_id" binding:"required"`
	Quantity               uint `json:"quantity" binding:"required"`
}
