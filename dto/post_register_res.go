package dto

import "seadeals-backend/model"

type RegisterResponse struct {
	ID       uint         `json:"id"`
	FullName string       `json:"full_name"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	Phone    string       `json:"phone"`
	Role     string       `json:"role"`
	Wallet   model.Wallet `json:"wallet"`
}
