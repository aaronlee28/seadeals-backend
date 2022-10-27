package dto

import "time"

type Receipt struct {
	SellerName    string             `json:"seller_name"`
	Buyer         BuyerReceipt       `json:"buyer"`
	OrderDetail   OrderDetailReceipt `json:"order_detail"`
	Transaction   TransactionReceipt `json:"transaction"`
	Courier       CourierReceipt     `json:"courier"`
	PaymentMethod string             `json:"payment_method"`
}

type BuyerReceipt struct {
	Name       string    `json:"name"`
	BoughtDate time.Time `json:"bought_date"`
	Address    string    `json:"address"`
}

type OrderDetailReceipt struct {
	TotalQuantity         uint                          `json:"total_quantity"`
	TotalOrder            float64                       `json:"total_order"`
	DeliveryPrice         float64                       `json:"delivery_price"`
	Total                 float64                       `json:"total"`
	GlobalVoucherForOrder *GlobalVoucherForOrderReceipt `json:"global_voucher_for_order"`
	ShopVoucher           *ShopVoucherReceipt           `json:"shop_voucher"`
	OrderItems            []*OrderItemReceipt           `json:"order_items"`
}

type GlobalVoucherForOrderReceipt struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	TotalReduce float64 `json:"total_reduce"`
}

type ShopVoucherReceipt struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	TotalReduce float64 `json:"total_reduce"`
}

type OrderItemReceipt struct {
	Name         string  `json:"name"`
	Weight       uint    `json:"weight"`
	Quantity     uint    `json:"quantity"`
	PricePerItem float64 `json:"price_per_item"`
	Discount     float64 `json:"discount"`
	Subtotal     float64 `json:"subtotal"`
	Variant      string  `json:"variant"`
}

type TransactionReceipt struct {
	TotalTransaction float64                  `json:"total_transaction"`
	GlobalDiscount   []*GlobalDiscountReceipt `json:"global_discount"`
	OrderPayments    []*OrderPaymentReceipt   `json:"order_payments"`
	Total            float64                  `json:"total"`
}

type GlobalDiscountReceipt struct {
	SellerName   string  `json:"seller_name"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	TotalReduced float64 `json:"total_reduced"`
}

type OrderPaymentReceipt struct {
	SellerName string  `json:"seller_name"`
	TotalOrder float64 `json:"total_order"`
}

type CourierReceipt struct {
	Name    string `json:"name"`
	Service string `json:"service"`
}
