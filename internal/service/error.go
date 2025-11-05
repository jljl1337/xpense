package service

import "fmt"

type ErrorCode int

const (
	ErrCodeBadRequest ErrorCode = iota
	ErrCodeUnauthorized
	ErrCodeForbidden
	ErrCodeNotFound
	ErrCodeConflict
	ErrCodeUnprocessable
	ErrCodeInternal
)

type ServiceError struct {
	Code    ErrorCode
	Message string
}

func NewServiceErrorf(code ErrorCode, format string, args ...any) *ServiceError {
	return NewServiceError(code, fmt.Sprintf(format, args...))
}

func NewServiceError(code ErrorCode, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}

func (e *ServiceError) Error() string {
	return e.Message
}
