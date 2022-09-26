package dto

type GoogleLogin struct {
	TokenID string `json:"token_id" binding:"required"`
}
