package dto

type RegisterSeaLabsPayReq struct {
	AccountNumber string `json:"account_number" binding:"required,numeric,len=16"`
	Name          string `json:"name" binding:"required"`
}
