package dto

type CreateAddressReq struct {
	SubDistrictID uint   `json:"sub_district_id" binding:"required"`
	Address       string `json:"address" binding:"required"`
	Zipcode       string `json:"zipcode" binding:"required"`
}
