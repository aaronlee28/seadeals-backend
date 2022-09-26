package dto

type UserJWT struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	WalletID uint   `json:"wallet_id"`
}

const (
	JWTAccessToken  = "access_token"
	JWTRefreshToken = "refresh_token"
)
