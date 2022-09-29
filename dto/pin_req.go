package dto

type PinReq struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}
