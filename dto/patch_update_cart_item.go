package dto

type UpdateCartItemReq struct {
	CartItemID      uint `json:"cart_item_id" binding:"required"`
	CurrentQuantity uint `json:"current_quantity" binding:"required"`
}
