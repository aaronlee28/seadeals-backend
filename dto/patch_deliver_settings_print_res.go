package dto

import "seadeals-backend/model"

type DeliverSettingsPrintRes struct {
	SellerID   uint `json:"seller_id"`
	AllowPrint bool `json:"allow_print"`
}

func (_ *DeliverSettingsPrintRes) DeliverySettingsFromSeller(t *model.Seller) *DeliverSettingsPrintRes {
	return &DeliverSettingsPrintRes{
		SellerID:   t.ID,
		AllowPrint: t.AllowPrint,
	}
}
