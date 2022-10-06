package dto

type PatchProductVariantDetail struct {
	Price         float64 `json:"price"`
	Variant1Value *string `json:"variant_1_value"`
	Variant2Value *string `json:"variant_2_value"`
	VariantCode   *string `json:"variant_code"`
	PictureURL    *string `json:"picture_url"`
	Stock         uint    `json:"stock"`
}
