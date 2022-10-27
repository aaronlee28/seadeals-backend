package dto

import "time"

type Thermal struct {
	Buyer          BuyerThermal            `json:"buyer"`
	Courier        CourierThermal          `json:"courier"`
	SellerName     string                  `json:"seller_name"`
	TotalWeight    uint                    `json:"total_weight"`
	Price          float64                 `json:"price"`
	DeliveryNumber string                  `json:"delivery_number"`
	OriginCity     string                  `json:"origin_city"`
	IssuedAt       time.Time               `json:"issued_at"`
	Products       []*ProductDetailThermal `json:"products"`
}

type BuyerThermal struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}

type CourierThermal struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type ProductDetailThermal struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Variant  string `json:"variant"`
	Quantity uint   `json:"quantity"`
}
