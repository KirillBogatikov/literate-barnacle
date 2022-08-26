package user

import (
	"literate-barnacle/service"
	"literate-barnacle/service/models"
	"literate-barnacle/service/validation"

	"github.com/google/uuid"
)

type LoginRequest struct {
	models.Credentials
}

func (l LoginRequest) Validate() validation.Map {
	result := make(validation.Map)

	if len(l.Login) == 0 {
		result["login"] = validation.StatusEmpty
	}
	if len(l.Password) == 0 {
		result["password"] = validation.StatusEmpty
	}

	return result
}

type LoginResponse struct {
	service.BaseResponse
	Token string `json:"token,omitempty"`
}

type SignUpRequest struct {
	User models.User `json:"user"`
}

func (s SignUpRequest) Validate() validation.Map {
	result := make(validation.Map)

	user := s.User.Validate()
	if !user.IsValid() {
		result["user"] = user
	}

	return result
}

type SignUpResponse struct {
	service.BaseResponse
	UserId uuid.UUID `json:"userId,omitempty"`
}

type Response struct {
	service.BaseResponse
	User *models.User `json:"user,omitempty"`
}

func userNotFoundResponse() service.BaseResponse {
	return service.BaseResponse{
		Error: "Пользователь не найден",
	}
}
