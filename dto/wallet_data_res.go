package dto

import "seadeals-backend/model"

type WalletDataRes struct {
	UserID       uint                 `json:"user_id"`
	Balance      float64              `json:"balance"`
	Status       *string              `json:"status"`
	Transactions []*model.Transaction `json:"transactions"`
}
