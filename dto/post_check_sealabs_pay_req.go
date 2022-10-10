package dto

type CheckSeaLabsPayReq struct {
	AccountNumber string `json:"account_number" binding:"required,numeric,len=16"`
}
