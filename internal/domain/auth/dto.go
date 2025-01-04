package auth

import "time"

type LogInDTO struct {
	Username string
	Password string
}

type TokenDTO struct {
	AccessToken  string
	RefreshToken string
}

type SignUpDTO struct {
	UserID   uint64
	Username string
	Password string
	Sex      string
	Birthday time.Time
}
