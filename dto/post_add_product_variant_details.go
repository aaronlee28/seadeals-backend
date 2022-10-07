package dto

type AddVariantAndDetails struct {
	Variant1Name          *string                 `json:"variant_1_name"`
	Variant2Name          *string                 `json:"variant_2_name"`
	ProductVariantDetails []*ProductVariantDetail `json:"product_variant_details"`
}
