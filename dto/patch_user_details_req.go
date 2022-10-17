package dto

import "time"

type ChangeUserDetails struct {
	Username  string    `json:"username,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	BirthDate time.Time `json:"birth_date,omitempty"`
}
