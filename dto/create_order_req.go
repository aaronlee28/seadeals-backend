package dto

type CreateOrderReq struct {
	OrderItemID []int `json:"order_item_id" binding:"required"`
}
