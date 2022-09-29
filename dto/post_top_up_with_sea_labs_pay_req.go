package dto

type TopUpWalletWithSeaLabsPayReq struct {
	Amount        float64 `json:"amount" binding:"required"`
	AccountNumber string  `json:"account_number" binding:"required,numeric,len=16"`
}
