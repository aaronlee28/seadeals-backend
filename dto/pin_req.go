package dto

type PinReq struct {
	Pin int `json:"pin" binding:"required"`
}
