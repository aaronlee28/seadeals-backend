package dto

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Username  string `json:"username" binding:"required,alphanum"`
	FullName  string `json:"full_name" binding:"required"`
	Gender    string `json:"gender" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	BirthDate string `json:"birth_date" binding:"required"`
}
