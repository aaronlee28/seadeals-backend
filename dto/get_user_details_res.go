package dto

import "time"

type UserDetailsRes struct {
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	AvatarURL string    `json:"avatar_url"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}
