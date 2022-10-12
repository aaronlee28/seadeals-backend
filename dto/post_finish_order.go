package dto

type FinishOrderReq struct {
	OrderID uint `json:"order_id" binding:"required"`
}
