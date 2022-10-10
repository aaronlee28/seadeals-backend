package dto

const (
	DeliveryWaitingForPayment = "waiting for payment"
	DeliveryWaitingForSeller  = "waiting for seller"
	DeliveryOngoing           = "on delivery"
	DeliveryFailed            = "failed"
	DeliveryDone              = "done"
)

const (
	OrderWaitingPayment = "waiting for payment"
	OrderWaitingSeller  = "waiting for seller"
	OrderOnDelivery     = "on delivery"
	OrderDelivered      = "delivered"
	OrderComplained     = "complained"
	OrderRefunded       = "refunded"
	OrderDone           = "done"
)
