package dto

type SeaDealspayReq struct {
	CardNumber string `json:"card_number" binding:"required,numeric"`
	Amount     int    `json:"amount" binding:"required"`
}
