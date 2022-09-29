package dto

type SeaLabsPayReq struct {
	Amount       string `json:"amount" binding:"required,numeric"`
	MerchantCode string `json:"merchant_code" binding:"required"`
	Message      string `json:"message" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
	Status       string `json:"status" binding:"required"`
	TxnID        string `json:"txn_id" binding:"required"`
}

const (
	TXN_PAID   = "TXN_PAID"
	TXN_FAILED = "TXN_FAILED"
)
