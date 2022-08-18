package service

import "regexp"

type ValidationError int

const (
	ValidationNoErr ValidationError = iota
	ErrValidationEmpty
	ErrValidationShort
	ErrValidationLong
	ErrValidationInvalid
)

type ValidationMap map[string]interface{}

func (v ValidationMap) IsValid() bool {
	return len(v) == 0
}

type Validatable interface {
	Validate() ValidationMap
}

func validateString(input string, min, max uint, regex *regexp.Regexp) ValidationError {
	length := uint(len(input))
	switch {
	case length == 0:
		return ErrValidationEmpty
	case length < min:
		return ErrValidationShort
	case length > max:
		return ErrValidationLong
	case regex != nil && !regex.MatchString(input):
		return ErrValidationInvalid
	default:
		return ValidationNoErr
	}
}
