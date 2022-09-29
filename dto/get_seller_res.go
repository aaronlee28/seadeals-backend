package dto

import (
	"seadeals-backend/model"
	"strconv"
	"time"
)

type GetSellerRes struct {
	ID            uint           `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Address       *GetAddressRes `json:"address"`
	ProfileURL    string         `json:"profile_url"`
	BannerURL     string         `json:"banner_url"`
	Followers     uint           `json:"followers"`
	Following     uint           `json:"following"`
	Rating        float64        `json:"rating"`
	TotalReviewer uint           `json:"total_reviewer"`
	JoinDate      string         `json:"join_date"`
	IsFollow      bool           `json:"is_follow"`
}

func (_ *GetSellerRes) From(s *model.Seller) *GetSellerRes {
	address := new(GetAddressRes).From(s.Address)

	joinStatus := "just now"
	now := time.Now()
	deltaYear := now.Year() - s.CreatedAt.Year()
	if deltaYear > 0 {
		plural := " years"
		if deltaYear == 1 {
			plural = " year"
		}
		joinStatus = strconv.Itoa(deltaYear) + plural + " ago"
	}

	deltaMonth := now.Month() - s.CreatedAt.Month()
	if deltaYear == 0 && deltaMonth > 0 {
		plural := " months"
		if deltaMonth == 1 {
			plural = " month"
		}
		joinStatus = strconv.Itoa(int(deltaMonth)) + plural + " ago"
	}

	deltaDay := now.Day() - s.CreatedAt.Day()
	if deltaYear == 0 && deltaMonth == 0 && deltaDay > 0 {
		plural := " days"
		if deltaDay == 1 {
			plural = " day"
		}
		joinStatus = strconv.Itoa(deltaDay) + plural + " ago"
	}

	isFollow := false
	if s.SocialGraph != nil {
		isFollow = true
	}

	return &GetSellerRes{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Address:     address,
		ProfileURL:  s.PictureURL,
		BannerURL:   s.BannerURL,
		JoinDate:    joinStatus,
		IsFollow:    isFollow,
	}
}
