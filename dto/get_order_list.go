package dto

import (
	"seadeals-backend/model"
	"time"
)

type OrderListRes struct {
	ID                       uint                  `json:"id"`
	BuyerName                string                `json:"buyer_name"`
	SellerID                 uint                  `json:"seller_id"`
	Seller                   SellerOrderList       `json:"seller"`
	VoucherID                uint                  `json:"voucher_id"`
	Voucher                  *VoucherOrderList     `json:"voucher"`
	TransactionID            uint                  `json:"transaction_id"`
	Transaction              TransactionOrderList  `json:"transaction"`
	TotalOrderPrice          float64               `json:"total_order_price"`
	TotalOrderPriceAfterDisc float64               `json:"total_order_price_after_disc"`
	TotalDelivery            float64               `json:"total_delivery"`
	Status                   string                `json:"status"`
	HasReviewedAllItem       bool                  `json:"has_reviewed_all_item"`
	OrderItems               []*OrderItemOrderList `json:"order_items"`
	DeliveryID               uint                  `json:"delivery_id"`
	Delivery                 *DeliveryOrderList    `json:"delivery"`
	Complaint                *model.Complaint      `json:"complaint"`
	UpdatedAt                time.Time             `json:"updated_at"`
}

type SellerOrderList struct {
	Name string `json:"name"`
}

type VoucherOrderList struct {
	Code          string  `json:"code"`
	VoucherType   string  `json:"voucher_type"`
	Amount        float64 `json:"amount"`
	AmountReduced float64 `json:"amount_reduced"`
}

type TransactionOrderList struct {
	PaymentMethod string     `json:"payment_method"`
	Total         float64    `json:"total"`
	Status        string     `json:"status"`
	PayedAt       *time.Time `json:"payed_at"`
}

type OrderItemOrderList struct {
	ID                     uint                   `json:"id"`
	ProductVariantDetailID uint                   `json:"product_variant_detail_id"`
	ProductDetail          ProductDetailOrderList `json:"product_detail"`
	Quantity               uint                   `json:"quantity"`
	Subtotal               float64                `json:"subtotal"`
}

type ProductDetailOrderList struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	CategoryID   uint             `json:"category_id"`
	Category     string           `json:"category"`
	Slug         string           `json:"slug"`
	PhotoURL     string           `json:"photo_url"`
	Variant      string           `json:"variant"`
	Price        float64          `json:"price"`
	ReviewByUser *ReviewOrderList `json:"review_by_user"`
}

type ReviewOrderList struct {
	ID          uint    `json:"id"`
	Rating      int     `json:"rating"`
	Description *string `json:"description"`
	ImageUrl    *string `json:"image_url"`
}

type DeliveryOrderList struct {
	DestinationAddress string                       `json:"destination_address"`
	Status             string                       `json:"status"`
	DeliveryNumber     string                       `json:"delivery_number"`
	ETA                int                          `json:"eta"`
	CourierID          uint                         `json:"courier_id"`
	Courier            string                       `json:"courier"`
	Activity           []*DeliveryActivityOrderList `json:"activity"`
}

type DeliveryActivityOrderList struct {
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
