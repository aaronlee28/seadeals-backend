package dto

type PatchVariantAndDetails struct {
	Variant1Name          *string                    `json:"variant_1_name"`
	Variant2Name          *string                    `json:"variant_2_name"`
	ProductVariantDetails *PatchProductVariantDetail `json:"product_variant_details"`
}
