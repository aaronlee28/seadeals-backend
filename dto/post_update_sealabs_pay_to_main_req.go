package dto

type UpdateSeaLabsPayToMainReq struct {
	UserID        uint   `json:"user_id" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required,numeric,len=16"`
}
