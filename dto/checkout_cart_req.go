package dto

type CheckoutCartReq struct {
	GlobalVoucherCode string          `json:"global_voucher_code"`
	Cart              []*CartPerStore `json:"cart_per_store" binding:"required"`
	PaymentMethod     string          `json:"payment_method" binding:"required"`
}
