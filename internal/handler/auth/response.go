package auth

import "time"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_Token"`
}

type CurrentUserResponse struct {
	UserID    uint64    `json:"user_id"`
	Username  string    `json:"username"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
}
