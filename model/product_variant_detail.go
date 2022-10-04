package model

import "gorm.io/gorm"

type ProductVariantDetail struct {
	gorm.Model      `json:"-"`
	ID              uint            `json:"id" gorm:"primaryKey"`
	ProductID       uint            `json:"product_id"`
	Product         *Product        `json:"product"`
	Price           float64         `json:"price"`
	Variant1Value   *string         `json:"variant1_value"`
	Variant2Value   *string         `json:"variant2_value"`
	Variant1ID      *uint           `json:"variant1_id"`
	Variant2ID      *uint           `json:"variant2_id"`
	VariantCode     *string         `json:"variant_code"`
	PictureURL      *string         `json:"picture_url"`
	Stock           uint            `json:"stock"`
	ProductVariant1 *ProductVariant `json:"product_variant1" gorm:"foreignKey:Variant1ID; references:ID"`
	ProductVariant2 *ProductVariant `json:"product_variant2" gorm:"foreignKey:Variant2ID; references:ID"`
}
