package ctx

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized = errors.New("authorization header required")
	ErrTokenExpired = errors.New("token expired")
	ErrForbidden    = errors.New("forbidden")
)

type ParseTokenError struct {
	origin error
}

func NewParseTokenError(origin error) ParseTokenError {
	return ParseTokenError{
		origin: origin,
	}
}

func (p ParseTokenError) Error() string {
	return fmt.Sprintf("invalid header: %v", p.origin)
}
