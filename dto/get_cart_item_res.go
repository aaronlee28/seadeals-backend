package dto

type CartItemRes struct {
	Quantity     int     `json:"quantity"`
	PricePerItem float64 `json:"price_per_item"`
	SellerID     uint    `json:"seller_id"`
	Subtotal     float64 `json:"subtotal"`
	ProductName  string  `json:"product_name"`
}