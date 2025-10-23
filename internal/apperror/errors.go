package apperror

import (
	"fmt"
	"net/http"
)

// FieldError represents validation error for a specific field
type FieldError struct {
	Field           string    `json:"field"`
	Code            ErrorCode `json:"code"`
	Message         string    `json:"message,omitempty"`
	Expect          string    `json:"expect,omitempty"`
	IsInternalError bool      `json:"-"`
}

// AppError represents the application error structure
type AppError struct {
	Code          ErrorCode    `json:"code"`
	HttpStatus    int          `json:"-"`
	Message       string       `json:"message"`
	NotFoundField string       `json:"not_found_field,omitempty"`
	FieldErrors   []FieldError `json:"field_errors,omitempty"`
	cause         error        `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s (code: %s, http: %d): %v", e.Message, e.Code, e.HttpStatus, e.cause)
	}
	return fmt.Sprintf("%s (code: %s, http: %d)", e.Message, e.Code, e.HttpStatus)
}

// Unwrap implements the Wrapper interface for error chaining
func (e *AppError) Unwrap() error {
	return e.cause
}

// WithCause adds an underlying cause to the error
func (e *AppError) WithCause(cause error) *AppError {
	e.cause = cause
	return e
}

// WithFieldErrors adds field validation errors
func (e *AppError) WithFieldErrors(fieldErrors []FieldError) *AppError {
	e.FieldErrors = fieldErrors
	return e
}

// HTTPStatus returns the HTTP status code for the error
func (e *AppError) HTTPStatus() int {
	return e.HttpStatus
}

// WithNotFoundField sets which field caused the not found error.
func (e *AppError) WithNotFoundField(field string) *AppError {
	e.NotFoundField = field
	return e
}

// New creates a new application error
func New(code ErrorCode, httpStatus int, message string) *AppError {
	return &AppError{
		Code:        code,
		HttpStatus:  httpStatus,
		Message:     message,
		FieldErrors: []FieldError{},
	}
}

// Wrap creates a new application error with an underlying cause
func Wrap(code ErrorCode, httpStatus int, message string, cause error) *AppError {
	return &AppError{
		Code:        code,
		HttpStatus:  httpStatus,
		Message:     message,
		cause:       cause,
		FieldErrors: []FieldError{},
	}
}

// NewFieldError creates a new field validation error
func NewFieldError(field string, code ErrorCode) FieldError {
	return FieldError{
		Field: field,
		Code:  code,
	}
}

// WithMessage sets a human-readable message for the field error and returns the updated FieldError.
func (f FieldError) WithMessage(msg string) FieldError {
	f.Message = msg
	return f
}

// WithExpect sets the expected valid value or format for the field error and returns the updated FieldError.
func (f FieldError) WithExpect(expect string) FieldError {
	f.Expect = expect
	return f
}

// Common Errors

// ValidationError creates a validation error with field-specific errors
func ValidationError(message string, fieldErrors []FieldError) *AppError {
	return New(ErrorCode("VALIDATION_ERROR"), http.StatusUnprocessableEntity, message).WithFieldErrors(fieldErrors)
}

// BadRequest
func BadRequest(message string, code ErrorCode) *AppError {
	return New(ErrorCode(code), http.StatusBadRequest, message)
}

// Unauthorized
func Unauthorized(message string, code ErrorCode) *AppError {
	return New(ErrorCode(code), http.StatusUnauthorized, message)
}

// Forbidden
func Forbidden(message string, code ErrorCode) *AppError {
	return New(ErrorCode(code), http.StatusForbidden, message)
}

// Conflict
func Conflict(message string, code ErrorCode) *AppError {
	return New(ErrorCode(code), http.StatusConflict, message)
}

// NotFound
func NotFound(message string, code ErrorCode) *AppError {
	return New(ErrorCode(code), http.StatusNotFound, message)
}

// InternalServer
func InternalServer(message string) *AppError {
	return New(ErrorCode("INTERNAL_ERROR"), http.StatusInternalServerError, message)
}

// Helper functions to work with AppError

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError returns the AppError if the error is one
func AsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}

// GetHttpStatus returns HTTP status code from error
// Returns 500 if error is not AppError
func GetHttpStatus(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPStatus()
	}
	return 500
}
