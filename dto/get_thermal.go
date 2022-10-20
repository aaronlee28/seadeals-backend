package dto

type Thermal struct {
	Buyer          BuyerThermal           `json:"buyer"`
	SellerName     string                 `json:"seller_name"`
	TotalWeight    uint                   `json:"total_weight"`
	DeliveryNumber string                 `json:"delivery_number"`
	OriginCity     string                 `json:"origin_city"`
	Product        []ProductDetailThermal `json:"product"`
}

type BuyerThermal struct {
	Name          string              `json:"name"`
	AddressDetail BuyerAddressThermal `json:"address_detail"`
}

type BuyerAddressThermal struct {
	Address     string `json:"address"`
	Province    string `json:"province"`
	City        string `json:"city"`
	SubDistrict string `json:"sub_district"`
	PostalCode  string `json:"postal_code"`
}

type ProductDetailThermal struct {
	Name     string `json:"name"`
	Variant  string `json:"variant"`
	Quantity uint   `json:"quantity"`
}
