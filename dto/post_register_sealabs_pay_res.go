package dto

import "seadeals-backend/model"

type RegisterSeaLabsPayRes struct {
	Status         string                       `json:"status"`
	SeaLabsAccount *model.UserSealabsPayAccount `json:"sea_labs_account"`
}
