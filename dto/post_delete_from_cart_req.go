package dto

type DeleteFromCartReq struct {
	OrderItemID uint `json:"order_item_id" binding:"required"`
	UserID      uint `json:"user_id" binding:"required"`
}
