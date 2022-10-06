package dto

type SellerCancelOrderReq struct {
	OrderID uint `json:"order_id" binding:"required"`
}
