package dto

type CheckoutCartReq struct {
	GlobalVoucherCode string          `json:"global_voucher_code" binding:"omitempty,alphanum"`
	Cart              []*CartPerStore `json:"cart_per_store" binding:"required"`
	PaymentMethod     string          `json:"payment_method" binding:"required"`
	AccountNumber     string          `json:"account_number"`
	BuyerAddressID    uint            `json:"buyer_address_id" binding:"required"`
}
