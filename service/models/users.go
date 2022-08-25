package models

import (
	"literate-barnacle/service/validation"
	"regexp"

	"github.com/google/uuid"
)

var (
	_validLoginRegex    = regexp.MustCompile(`[a-zA-Zа-яА-я0-9_\-+*!@#$%^]+`)
	_validPasswordRegex = regexp.MustCompile(`[^ ]+`)
)

type Role uint8

const (
	RoleUnknown Role = iota
	RoleUser
	RoleAdmin
)

type Credentials struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (c Credentials) Validate() validation.Map {
	result := make(validation.Map)

	vErr := validation.ValidateString(c.Login, 4, 16, _validLoginRegex)
	if vErr != validation.StatusOk {
		result["login"] = vErr
	}

	vErr = validation.ValidateString(c.Password, 8, 128, _validPasswordRegex)
	if vErr != validation.StatusOk {
		result["password"] = vErr
	}

	return result
}

type User struct {
	Id          uuid.UUID   `json:"id,omitempty"`
	Credentials Credentials `json:"credentials,omitempty"`
	Role        Role        `json:"role,omitempty"`
}

func (u User) Validate() validation.Map {
	result := make(validation.Map)

	credentials := u.Credentials.Validate()
	if !credentials.IsValid() {
		result["credentials"] = credentials
	}

	return result
}
