package dto

type UpdateAddressRes struct {
	ID            uint   `json:"id""`
	Address       string `json:"address"`
	Zipcode       string `json:"zipcode"`
	SubDistrictID uint   `json:"sub_district_id"`
	UserID        uint   `json:"user_id"`
}
