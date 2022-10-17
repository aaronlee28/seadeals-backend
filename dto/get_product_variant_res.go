package dto

import (
	"seadeals-backend/model"
)

type ProductVariantRes struct {
	MinPrice        float64                 `json:"min_price"`
	MaxPrice        float64                 `json:"max_price"`
	ProductVariants []*GetProductVariantRes `json:"product_variants"`
}

type GetProductVariantRes struct {
	ID            uint    `json:"id"`
	ProductID     uint    `json:"product_id"`
	Price         float64 `json:"price"`
	Variant1Name  string  `json:"variant1_name"`
	Variant2Name  string  `json:"variant2_name"`
	Variant1Value *string `json:"variant1_value"`
	Variant2Value *string `json:"variant2_value"`
	VariantCode   *string `json:"variant_code"`
	PictureURL    *string `json:"picture_url"`
	Stock         uint    `json:"stock"`
}

func (_ *GetProductVariantRes) From(pv *model.ProductVariantDetail) *GetProductVariantRes {
	var name1, name2 string
	if pv.ProductVariant1 != nil {
		name1 = pv.ProductVariant1.Name
	}
	if pv.ProductVariant2 != nil {
		name2 = pv.ProductVariant2.Name
	}

	var variant1Value, variant2Value, pictureURL *string
	if pv.Variant1Value != nil {
		variant1Value = pv.Variant1Value
	}
	if pv.Variant2Value != nil {
		variant2Value = pv.Variant2Value
	}
	if pv.PictureURL != nil {
		pictureURL = pv.PictureURL
	}

	return &GetProductVariantRes{
		ID:            pv.ID,
		ProductID:     pv.ProductID,
		Price:         pv.Price,
		Variant1Name:  name1,
		Variant2Name:  name2,
		Variant1Value: variant1Value,
		Variant2Value: variant2Value,
		VariantCode:   pv.VariantCode,
		PictureURL:    pictureURL,
		Stock:         pv.Stock,
	}
}
