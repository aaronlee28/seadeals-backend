package dto

import "seadeals-backend/model"

type GetProductVariantRes struct {
	ID            uint    `json:"id"`
	ProductID     uint    `json:"product_id"`
	Price         float64 `json:"price"`
	Variant1Name  string  `json:"variant1_name"`
	Variant2Name  string  `json:"variant2_name"`
	Variant1Value string  `json:"variant1_value"`
	Variant2Value string  `json:"variant2_value"`
	VariantCode   string  `json:"variant_code"`
	PictureURL    string  `json:"picture_url"`
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
	return &GetProductVariantRes{
		ID:            pv.ID,
		ProductID:     pv.ProductID,
		Price:         pv.Price,
		Variant1Name:  name1,
		Variant2Name:  name2,
		Variant1Value: pv.Variant1Value,
		Variant2Value: pv.Variant2Value,
		VariantCode:   pv.VariantCode,
		PictureURL:    pv.PictureURL,
		Stock:         pv.Stock,
	}
}
