package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct{}

func NewBCrypt() BCrypt {
	return BCrypt{}
}

func (b BCrypt) Encrypt(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashed), err
}

func (b BCrypt) Compare(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return ErrMismatched
	case err != nil:
		return err
	default:
		return nil
	}
}
