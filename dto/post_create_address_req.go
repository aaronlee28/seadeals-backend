package dto

type CreateAddressReq struct {
	CityID      string `json:"city_id" binding:"required,numeric"`
	ProvinceID  string `json:"province_id" binding:"required,numeric"`
	Province    string `json:"province" binding:"required"`
	City        string `json:"city" binding:"required"`
	Type        string `json:"type" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required,numeric"`
	SubDistrict string `json:"sub_district" binding:"required"`
	Address     string `json:"address" binding:"required"`
}
