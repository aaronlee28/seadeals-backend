package dto

type UpdateAddressReq struct {
	UserID        uint   `json:"user_id" binding:"required"`
	SubDistrictID uint   `json:"sub_district_id" binding:"required"`
	Address       string `json:"address" binding:"required"`
	Zipcode       string `json:"zipcode" binding:"required"`
	ID            uint   `json:"id" binding:"required"`
}
