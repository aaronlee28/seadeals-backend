package dto

import "seadeals-backend/model"

type GetSellerRes struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     *GetAddressRes `json:"address"`
	PictureURL  string         `json:"picture_url"`
	BannerURL   string         `json:"banner_url"`
}

func (_ *GetSellerRes) From(s *model.Seller) *GetSellerRes {
	address := new(GetAddressRes).From(s.Address)
	return &GetSellerRes{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Address:     address,
		PictureURL:  s.PictureURL,
		BannerURL:   s.BannerURL,
	}
}
