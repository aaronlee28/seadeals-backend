package dto

type TotalPredictedPriceRes struct {
	PredictedPrices     []*PredictedPriceRes `json:"predicted_prices"`
	GlobalVoucherID     *uint                `json:"global_voucher_id"`
	TotalPredictedPrice float64              `json:"total_predicted_price"`
}
