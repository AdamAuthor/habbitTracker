package models

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
