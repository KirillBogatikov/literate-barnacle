package hash

import (
	"errors"
)

type Encryptor interface {
	Encrypt(input string) (string, error)
	Compare(input, hash string) error
}

var (
	ErrMismatched = errors.New("hash mismatch")
)
