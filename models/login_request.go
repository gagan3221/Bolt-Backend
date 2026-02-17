package models

type LoginRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}
