package dto

type DeliveryCalculateReq struct {
	OriginCity      string `json:"origin_city"`
	DestinationCity string `json:"destination_city"`
	Weight          string `json:"weight"`
	Courier         string `json:"courier"`
}

type DeliveryCalculateRes struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Costs []struct {
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost        []struct {
			Value int    `json:"value"`
			Etd   string `json:"etd"`
			Note  string `json:"note"`
		} `json:"cost"`
	} `json:"costs"`
}

type DeliveryCalculateReturn struct {
	Total int `json:"total"`
	Eta   int `json:"eta"`
}
