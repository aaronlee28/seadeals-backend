package dto

type AddDeliveryReq struct {
	CourierID  uint  `json:"courier_id" binding:"required"`
	IsSelected *bool `json:"is_selected" binding:"required"`
	SlaDay     int   `json:"sla_day"`
}
