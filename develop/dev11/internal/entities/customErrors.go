package entities

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidData   = errors.New("invalid data")
	ErrISE           = errors.New("internal server error")
)
