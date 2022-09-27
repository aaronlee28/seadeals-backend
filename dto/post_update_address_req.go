package dto

type UpdateAddressReq struct {
	ID          uint   `json:"id" binding:"required"`
	UserID      uint   `json:"user_id" binding:"required"`
	CityID      string `json:"city_id"`
	ProvinceID  string `json:"province_id"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Type        string `json:"type"`
	PostalCode  string `json:"postal_code"`
	SubDistrict string `json:"sub_district"`
	Address     string `json:"address"`
}
