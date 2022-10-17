package dto

type UserJWT struct {
	Name      string  `json:"name"`
	UserID    uint    `json:"user_id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	WalletID  uint    `json:"wallet_id"`
	AvatarURL *string `json:"avatar_url"`
}

const (
	JWTAccessToken  = "access_token"
	JWTRefreshToken = "refresh_token"
)
