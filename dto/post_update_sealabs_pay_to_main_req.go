package dto

type UpdateSeaLabsPayToMainReq struct {
	AccountNumber string `json:"account_number" binding:"required,numeric,len=16"`
}
