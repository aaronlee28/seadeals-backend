package dto

type CartItemRes struct {
	ID           uint    `json:"id"`
	Quantity     uint    `json:"quantity"`
	PricePerItem float64 `json:"price_per_item"`
	SellerID     uint    `json:"seller_id"`
	SellerName   string  `json:"seller_name"`
	ImageURL     string  `json:"image_url"`
	Subtotal     float64 `json:"subtotal"`
	ProductName  string  `json:"product_name"`
}
