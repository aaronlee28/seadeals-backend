package dto

type SellerCancelOrderReq struct {
	OrderID      uint   `json:"order_id" binding:"required"`
	OrderItemsID []uint `json:"order_items_id"`
}
