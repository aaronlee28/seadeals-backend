package dto

type RegisterSeaLabsPayReq struct {
	UserID        uint   `json:"user_id" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required,numeric,len=16"`
	Name          string `json:"name" binding:"required"`
}
