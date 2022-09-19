package dto

type GoogleLogin struct {
	Email string `json:"email" binding:"required"`
}
