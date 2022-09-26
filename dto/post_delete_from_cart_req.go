package dto

type DeleteFromCartReq struct {
	CartItemID uint `json:"cart_item_id" binding:"required"`
	UserID     uint `json:"user_id" binding:"required"`
}
