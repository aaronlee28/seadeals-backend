package dto

type ChangePasswordReq struct {
	Email             string `json:"email"`
	CurrentPassword   string `json:"current_password"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
}
