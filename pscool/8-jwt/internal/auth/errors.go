package auth

import "errors"

var (
	ErrInvalidSessionId = errors.New("invalid session id")
	ErrInvalidCode      = errors.New("invalid code")
)
