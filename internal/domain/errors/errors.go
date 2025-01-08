package errors

import "errors"

var (
	// ErrMissingCredentials AuthService
	ErrMissingCredentials = errors.New("missing credentials")

	// ErrUserAlreadyExists AuthService
	ErrUserAlreadyExists = errors.New("username already exists")

	// ErrUserNotValid AuthService
	ErrUserNotValid = errors.New("username/password is not correct")

	// ErrSessionNotFound AuthService
	ErrSessionNotFound = errors.New("session not found")

	// ErrRefreshTokenExpired AuthService
	ErrRefreshTokenExpired = errors.New("refresh token expired")

	// ErrUserNotFound UserService
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCode GroupService
	ErrInvalidCode = errors.New("invalid code")

	// ErrAlreadyMember GroupService
	ErrAlreadyMember = errors.New("already a member")

	// ErrGroupNotFound GroupService
	ErrGroupNotFound = errors.New("group not found")

	// ErrMemberNotFound GroupService
	ErrMemberNotFound = errors.New("member not found")

	// ErrAreNotOwner GroupService
	ErrAreNotOwner = errors.New("you aren't owner")

	// ErrOwnerCannotLeave GroupService
	ErrOwnerCannotLeave = errors.New("owner cannot leave")

	// ErrCannotKickYourself GroupService
	ErrCannotKickYourself = errors.New("you can not kick yourself")

	// ErrProductNotFound ProductService
	ErrProductNotFound = errors.New("product not found")
)
