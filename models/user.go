package models

type UserDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	OauthSub string `json:"id"`
}

