package dto

type PostCreateProductReq struct {
	//used to create product
	Name          string `json:"name" binding:"required"`
	CategoryID    uint   `json:"category_id" binding:"required"`
	IsBulkEnabled bool   `json:"is_bulk_enabled"`
	MinQuantity   uint   `json:"min_quantity"`
	MaxQuantity   uint   `json:"max_quantity"`
	//used to create product detail
	ProductDetail *ProductDetailReq `json:"product_detail_req" binding:"required"`
	//used to create product_photos
	ProductPhotos []*ProductPhoto `json:"product_photos" binding:"required"`
	//used to create product_variant
	HasVariant            bool                  `json:"has_variant" binding:"required"`
	Variant1Name          string                `json:"variant_1_name"`
	Variant2Name          *string               `json:"variant_2_name"`
	ProductVariantDetails *ProductVariantDetail `json:"product_variant_details" binding:"required"`
}
