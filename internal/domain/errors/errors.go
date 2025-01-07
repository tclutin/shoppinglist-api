package errors

import "errors"

var (
	ErrMissingCredentials  = errors.New("missing credentials")
	ErrUserAlreadyExists   = errors.New("username already exists")
	ErrUserNotValid        = errors.New("username/password is not correct")
	ErrSessionNotFound     = errors.New("session not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCode        = errors.New("invalid code")
	ErrAlreadyMember      = errors.New("already a member")
	ErrGroupNotFound      = errors.New("group not found")
	ErrMemberNotFound     = errors.New("member not found")
	ErrAreNotOwner        = errors.New("you aren't owner")
	ErrOwnerCannotLeave   = errors.New("owner cannot leave")
	ErrCannotKickYourself = errors.New("you can not kick yourself")
)
