package dto

type TopUpWalletWithSeaLabsPayReq struct {
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"account_number"`
}
