package domainerrors

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTooManyAttempts    = errors.New("invalid credentials")
)
