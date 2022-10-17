package dto

type WalletTransactionsQuery struct {
	SortBy string `json:"sortBy"`
	Sort   string `json:"sort"`
	Page   string `json:"page"`
	Limit  string `json:"limit"`
}

const (
	SeaLabsPay = "sealabs pay"
	Wallet     = "wallet"
)
