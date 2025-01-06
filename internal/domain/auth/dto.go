package auth

import "github.com/google/uuid"

type LogInDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

type TokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDTO struct {
	RefreshToken uuid.UUID
}
