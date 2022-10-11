package dto

type DeleteFromCartReq struct {
	CartItemID uint `json:"cart_item_id" binding:"required"`
}
