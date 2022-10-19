package dto

type PostCreateProductReq struct {

	//used to create product
	Name          string `json:"name" binding:"required"`
	CategoryID    uint   `json:"category_id" binding:"required"`
	IsBulkEnabled bool   `json:"is_bulk_enabled"`
	MinQuantity   uint   `json:"min_quantity"`
	MaxQuantity   uint   `json:"max_quantity"`
	//used to create product detail
	ProductDetail *ProductDetailsReq `json:"product_detail_req" binding:"required"`
	//used to create product_photos
	ProductPhotos []*ProductPhoto `json:"product_photos" binding:"required"`
	//used to create product_variant
	DefaultPrice *float64             `json:"default_price"`
	DefaultStock *uint                `json:"default_stock"`
	Variant1Name *string              `json:"variant_1_name"`
	Variant2Name *string              `json:"variant_2_name"`
	VariantArray []*VariantAndDetails `json:"variant_array"`
}
