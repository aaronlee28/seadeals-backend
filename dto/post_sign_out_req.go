package dto

type SignOutReq struct {
	UserID uint `json:"user_id" binding:"required"`
}
