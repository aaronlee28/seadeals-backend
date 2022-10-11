package dto

type CartPerStore struct {
	VoucherCode string `json:"voucher_code"`
	SellerID    uint   `json:"seller_id"`
	CourierID   uint   `json:"courier_id" binding:"required"`
	CartItemID  []uint `json:"cart_item_id"`
}
