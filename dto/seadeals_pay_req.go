package dto

type SeaDealspayReq struct {
	CardNumber string `json:"card_number" binding:"required"`
	Amount     int    `json:"amount" binding:"required"`
}
