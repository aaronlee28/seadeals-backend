package dto

type PayWithSeaLabsPayReq struct {
	Amount        int    `json:"amount" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required,numeric,len=16"`
}
