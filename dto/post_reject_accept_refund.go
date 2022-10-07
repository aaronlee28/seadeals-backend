package dto

import "seadeals-backend/model"

type RejectAcceptRefundReq struct {
	OrderID uint `json:"order_id" binding:"required"`
}

type RejectAcceptRefundRes struct {
	Order          *model.Order `json:"order"`
	AmountRefunded float64      `json:"amount_refunded"`
}
