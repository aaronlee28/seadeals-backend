package dto

import "seadeals-backend/model"

type ProductVariantDetail struct {
	Price         float64 `json:"price" binding:"required"`
	Variant1Value *string `json:"variant_1_value" binding:"required"`
	Variant2Value *string `json:"variant_2_value"`
	VariantCode   *string `json:"variant_code"`
	PictureURL    *string `json:"picture_url"`
	Stock         uint    `json:"stock" binding:"required"`
}

func (_ *ProductVariantDetail) From(pvd *model.ProductVariantDetail) *ProductVariantDetail {
	return &ProductVariantDetail{
		Price:         pvd.Price,
		Variant1Value: pvd.Variant1Value,
		Variant2Value: pvd.Variant2Value,
		VariantCode:   pvd.VariantCode,
		PictureURL:    pvd.PictureURL,
		Stock:         pvd.Stock,
	}
}
