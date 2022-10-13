package dto

type CartItemRes struct {
	ID                  uint     `json:"id"`
	Quantity            uint     `json:"quantity"`
	ProductVariant      string   `json:"product_variant"`
	MinQuantity         uint     `json:"min_quantity"`
	MaxQuantity         uint     `json:"max_quantity"`
	Stock               uint     `json:"stock"`
	ProductSlug         string   `json:"product_slug"`
	DiscountPercent     *int     `json:"discount_percent"`
	DiscountNominal     *float64 `json:"discount_nominal"`
	PriceBeforeDiscount float64  `json:"price_before_discount"`
	PricePerItem        float64  `json:"price_per_item"`
	SellerID            uint     `json:"seller_id"`
	SellerName          string   `json:"seller_name"`
	ImageURL            string   `json:"image_url"`
	Subtotal            float64  `json:"subtotal"`
	ProductName         string   `json:"product_name"`
}
