package service

import "literate-barnacle/service/validation"

type BaseResponse struct {
	Error      string         `json:"error,omitempty"`
	Validation validation.Map `json:"validation,omitempty"`
}

func (b BaseResponse) IsSuccess() bool {
	return len(b.Error) == 0 && b.Validation.IsValid()
}
