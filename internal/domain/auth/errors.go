package auth

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordTooLong    = errors.New("password is too long")
	ErrLoginTooLong       = errors.New("login is too long")
)
