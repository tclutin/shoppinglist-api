package errors

import "errors"

var (
	ErrMissingCredentials  = errors.New("missing credentials")
	ErrUserAlreadyExists   = errors.New("username already exists")
	ErrUserNotValid        = errors.New("username/password is not correct")
	ErrSessionNotFound     = errors.New("session not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")

	ErrUserNotFound = errors.New("user not found")
)
