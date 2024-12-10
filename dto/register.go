package dto

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
