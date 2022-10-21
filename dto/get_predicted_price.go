package dto

type PredictedPriceReq struct {
	GlobalVoucherCode string          `json:"global_voucher_code" binding:"omitempty,alphanum"`
	Cart              []*CartPerStore `json:"cart_per_store" binding:"required"`
	BuyerAddressID    uint            `json:"buyer_address_id" binding:"required"`
}

type PredictedPriceRes struct {
	SellerID       uint    `json:"seller_id"`
	VoucherID      *uint   `json:"voucher_id"`
	TotalOrder     float64 `json:"total_order"`
	DeliveryPrice  float64 `json:"delivery_price"`
	PredictedPrice float64 `json:"predicted_price"`
}
