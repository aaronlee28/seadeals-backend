package dto

type StepUpPasswordRes struct {
	Password string `json:"password" binding:"required"`
}
