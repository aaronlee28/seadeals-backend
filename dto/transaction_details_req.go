package dto

type TransactionDetailsReq struct {
	TransactionID uint `json:"transaction_id" binding:"required"`
}
