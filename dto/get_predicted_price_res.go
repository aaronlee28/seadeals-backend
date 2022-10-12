package dto

type PredictedPriceRes struct {
	SellerID       uint    `json:"seller_id"`
	VoucherID      *uint   `json:"voucher_id"`
	TotalOrder     float64 `json:"total_order"`
	DeliveryPrice  float64 `json:"delivery_price"`
	PredictedPrice float64 `json:"predicted_price"`
}
