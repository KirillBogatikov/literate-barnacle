package validation

import "regexp"

type Status int

const (
	StatusOk Status = iota
	StatusEmpty
	StatusShort
	StatusLong
	StatusInvalid
)

type Map map[string]interface{}

func (v Map) IsValid() bool {
	return len(v) == 0
}

type Validatable interface {
	Validate() Map
}

func ValidateString(input string, min, max uint, regex *regexp.Regexp) Status {
	length := uint(len(input))
	switch {
	case length == 0:
		return StatusEmpty
	case length < min:
		return StatusShort
	case length > max:
		return StatusLong
	case regex != nil && !regex.MatchString(input):
		return StatusInvalid
	default:
		return StatusOk
	}
}
