package dto

import "seadeals-backend/model"

type GetAddressRes struct {
	ID          uint   `json:"id"`
	CityID      string `json:"city_id"`
	ProvinceID  string `json:"province_id"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Type        string `json:"type"`
	PostalCode  string `json:"postal_code"`
	SubDistrict string `json:"sub_district"`
	Address     string `json:"address"`
	IsMain      bool   `json:"is_main"`
}

func (_ *GetAddressRes) From(address *model.Address) *GetAddressRes {
	res := &GetAddressRes{
		ID:          address.ID,
		CityID:      address.CityID,
		ProvinceID:  address.ProvinceID,
		Province:    address.Province,
		City:        address.City,
		Type:        address.Type,
		PostalCode:  address.PostalCode,
		SubDistrict: address.SubDistrict,
		Address:     address.Address,
		IsMain:      address.IsMain,
	}

	return res
}
