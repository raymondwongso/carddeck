package entity

const (
	ErrParamInvalid    = "common.parameter_invalid"
	ErrMsgParamInvalid = "invalid parameter"
	ErrInternal        = "common.internal"
	ErrMsgInternal     = "something wrong happened"
)

type errorStruct struct {
	Code    string         `json:"code"`
	Message string         ` json:"message"`
	Details []*errorDetail `json:"error_details"`
}

type errorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewError creates new erro
func NewError(code, message string) *errorStruct {
	return &errorStruct{
		Code:    code,
		Message: message,
	}
}

// NewErrorDetail creates new error detail
func NewErrorDetail(field, message string) *errorDetail {
	return &errorDetail{
		Field:   field,
		Message: message,
	}
}

// AddDetail add details to errorResponse
func (e *errorStruct) AddDetail(errorDetail *errorDetail) {
	e.Details = append(e.Details, errorDetail)
}
