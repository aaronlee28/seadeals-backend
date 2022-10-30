package dto

import (
	"seadeals-backend/model"
	"strconv"
	"time"
)

type GetSellerRes struct {
	ID            uint           `json:"id"`
	UserID        uint           `json:"user_id"`
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
	TotalProduct  int64          `json:"total_product"`
}

func (_ *GetSellerRes) From(s *model.Seller) *GetSellerRes {
	address := new(GetAddressRes).From(s.Address)

	joinStatus := "just now"
	now := time.Now()
	deltaTime := now.Sub(s.CreatedAt).Hours()

	day := deltaTime / 24
	if int(day) > 0 {
		joinStatus = "joined " + strconv.Itoa(int(day))
		if day == 1 {
			joinStatus += " day ago"
		} else {
			joinStatus += " days ago"
		}
	}

	month := day / 30
	if int(month) > 0 {
		joinStatus = "joined " + strconv.Itoa(int(month))
		if month == 1 {
			joinStatus += " month ago"
		} else {
			joinStatus += " months ago"
		}
	}

	year := month / 12
	if int(year) > 0 {
		joinStatus = "joined " + strconv.Itoa(int(year))
		if year == 1 {
			joinStatus += " year ago"
		} else {
			joinStatus += " years ago"
		}
	}

	isFollow := false
	if s.SocialGraph != nil {
		isFollow = true
	}

	avatarURL := "https://firebasestorage.googleapis.com/v0/b/bucket-seadeals.appspot.com/o/avatars%2Fuser%2Fanonym.jpeg?alt=media&token=66dbb36a-2ac1-4b1f-ad67-b2834eefdcef"
	if s.User.AvatarURL != nil {
		avatarURL = *s.User.AvatarURL
	}

	return &GetSellerRes{
		ID:          s.ID,
		UserID:      s.UserID,
		Name:        s.Name,
		Description: s.Description,
		Address:     address,
		ProfileURL:  avatarURL,
		BannerURL:   s.BannerURL,
		JoinDate:    joinStatus,
		IsFollow:    isFollow,
	}
}
