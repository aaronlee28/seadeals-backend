package dto

type PaginatedTransactionsRes struct {
	TotalLength  int               `json:"totalLength"`
	Transactions []TransactionsRes `json:"transactions"`
}
