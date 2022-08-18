package service

import (
	"regexp"

	jwt "github.com/Viva-Victoria/bear-jwt"
	"github.com/google/uuid"
)

var (
	_validLoginRegex    = regexp.MustCompile(`[a-zA-Zа-яА-я0-9_\-+*!@#$%^]+`)
	_validPasswordRegex = regexp.MustCompile(`[^ ]+`)
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() ValidationMap {
	result := make(ValidationMap)

	if len(l.Login) == 0 {
		result["login"] = ErrValidationEmpty
	}
	if len(l.Password) == 0 {
		result["password"] = ErrValidationEmpty
	}

	return result
}

type LoginResponse struct {
	BaseResponse
	Token string `json:"token,omitempty"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c Credentials) Validate() ValidationMap {
	result := make(ValidationMap)

	vErr := validateString(c.Login, 4, 16, _validLoginRegex)
	if vErr != ValidationNoErr {
		result["login"] = vErr
	}

	vErr = validateString(c.Password, 8, 128, _validPasswordRegex)
	if vErr != ValidationNoErr {
		result["password"] = vErr
	}

	return result
}

type User struct {
	Id          uuid.UUID   `json:"id,omitempty"`
	Credentials Credentials `json:"credentials"`
}

func (u User) Validate() ValidationMap {
	result := make(ValidationMap)

	credentials := u.Credentials.Validate()
	if !credentials.IsValid() {
		result["credentials"] = credentials
	}

	return result
}

type SignUpRequest struct {
	User User `json:"user"`
}

func (s SignUpRequest) Validate() ValidationMap {
	result := make(ValidationMap)

	user := s.User.Validate()
	if !user.IsValid() {
		result["user"] = user
	}

	return result
}

type SignUpResponse struct {
	BaseResponse
	UserId uuid.UUID `json:"userId,omitempty"`
}

type TokenClaims struct {
	jwt.BasicClaims
	UserId string `json:"uid"`
}
