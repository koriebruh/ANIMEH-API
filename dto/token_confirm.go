package dto

type TokenConfirm struct {
	Token string `json:"token" binding:"required"`
}
