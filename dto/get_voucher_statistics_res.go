package dto

type VoucherStatisticsRes struct {
	Sale      float64 `json:"sale"`
	Order     uint    `json:"order"`
	UsageRate float64 `json:"usage_rate"`
	Buyer     uint    `json:"buyer"`
}
