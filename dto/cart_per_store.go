package dto

type CartPerStore struct {
	VoucherCode string `json:"voucher_code"`
	SellerID    uint   `json:"seller_id"`
	CartItemID  []uint `json:"cart_item_id"`
}
