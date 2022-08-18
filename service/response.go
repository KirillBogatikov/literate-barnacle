package service

type BaseResponse struct {
	Error      string        `json:"error"`
	Validation ValidationMap `json:"validation"`
}

func (b BaseResponse) IsSuccess() bool {
	return len(b.Error) == 0 && b.Validation.IsValid()
}
