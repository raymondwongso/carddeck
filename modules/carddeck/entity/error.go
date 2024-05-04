package entity

import "fmt"

const (
	ErrParamInvalid    = "common.parameter_invalid"
	ErrMsgParamInvalid = "invalid parameter"
	ErrInternal        = "common.internal"
	ErrMsgInternal     = "something wrong happened"

	ErrCardCodeInvalid    = "carddeck.card.code_invalid"
	ErrMsgCardCodeInvalid = "unknown card code"
)

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details []*ErrorDetail `json:"error_details"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements error interface
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates new erro
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorDetail creates new error detail
func NewErrorDetail(field, message string) *ErrorDetail {
	return &ErrorDetail{
		Field:   field,
		Message: message,
	}
}

// AddDetail add details to errorResponse
func (e *Error) AddDetail(errorDetail *ErrorDetail) {
	e.Details = append(e.Details, errorDetail)
}
