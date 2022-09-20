package dto

type CartItemRes struct {
	Quantity     int     `json:"quantity"`
	PricePerItem float64 `json:"price_per_item"`
	Subtotal     float64 `json:"subtotal"`
	ProductName  string  `json:"product_name"`
}
