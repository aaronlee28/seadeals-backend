package dto

import "seadeals-backend/model"

type GetAddressRes struct {
	ID          uint   `json:"id"`
	Address     string `json:"address"`
	Zipcode     string `json:"zipcode"`
	SubDistrict string `json:"sub_district"`
	District    string `json:"district"`
	City        string `json:"city"`
	Province    string `json:"province"`
}

func (_ *GetAddressRes) From(address *model.Address) *GetAddressRes {
	return &GetAddressRes{
		ID:          address.ID,
		Address:     address.Address,
		Zipcode:     address.Zipcode,
		SubDistrict: address.SubDistrict.Name,
		District:    address.SubDistrict.District.Name,
		City:        address.SubDistrict.District.City.Name,
		Province:    address.SubDistrict.District.City.Province.Name,
	}
}
