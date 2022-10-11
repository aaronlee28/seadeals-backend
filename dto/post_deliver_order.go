package dto

type DeliverOrderReq struct {
	OrderID uint `json:"order_id" binding:"required"`
}
