package dto

type UpdateAddressReq struct {
	ID          uint   `json:"id" binding:"required"`
	CityID      string `json:"city_id" binding:"omitempty,numeric"`
	ProvinceID  string `json:"province_id" binding:"omitempty,numeric"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Type        string `json:"type"`
	PostalCode  string `json:"postal_code"`
	SubDistrict string `json:"sub_district"`
	Address     string `json:"address"`
}
