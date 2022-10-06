package dto

type ProductVariantDetail struct {
	Price         float64 `json:"price" binding:"required"`
	Variant1Value *string `json:"variant_1_value" binding:"required"`
	Variant2Value *string `json:"variant_2_value"`
	VariantCode   *string `json:"variant_code"`
	PictureURL    *string `json:"picture_url"`
	Stock         uint    `json:"stock" binding:"required"`
}
