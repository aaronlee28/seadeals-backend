package dto

type PaginatedTransactionsRes struct {
	TotalLength  int               `json:"total_length"`
	TotalPage    int               `json:"total_page"`
	CurrentPage  int               `json:"current_page"`
	Limit        int               `json:"limit"`
	Transactions []TransactionsRes `json:"transactions"`
}
