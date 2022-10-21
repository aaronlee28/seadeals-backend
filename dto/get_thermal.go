package dto

type Thermal struct {
	Buyer          BuyerThermal            `json:"buyer"`
	SellerName     string                  `json:"seller_name"`
	TotalWeight    uint                    `json:"total_weight"`
	DeliveryNumber string                  `json:"delivery_number"`
	OriginCity     string                  `json:"origin_city"`
	Products       []*ProductDetailThermal `json:"product"`
}

type BuyerThermal struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}

type ProductDetailThermal struct {
	Name     string `json:"name"`
	Variant  string `json:"variant"`
	Quantity uint   `json:"quantity"`
}
