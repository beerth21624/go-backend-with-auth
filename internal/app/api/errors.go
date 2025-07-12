package api

import (
	"fmt"
	"net/http"

	"beerdosan-backend/internal/pkg/validator"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
	Details    map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code string, httpStatus int, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
		Err:        err,
	}
}

func NewInternalError(err error) *AppError {
	return NewAppError("INTERNAL_ERROR", http.StatusInternalServerError, "An unexpected error occurred", err)
}

func NewValidationError(err *validator.ValidationError) *AppError {
	details := make(map[string]interface{})

	for _, fieldErr := range err.Fields() {
		details[fieldErr.Field()] = map[string]string{
			"code":   fieldErr.Code(),
			"reason": fieldErr.Reason(),
		}
	}

	return &AppError{
		Code:       "VALIDATION_ERROR",
		HTTPStatus: http.StatusBadRequest,
		Message:    "Validation failed",
		Err:        err,
		Details:    details,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError("UNAUTHORIZED", http.StatusUnauthorized, message, nil)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError("FORBIDDEN", http.StatusForbidden, message, nil)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError("BAD_REQUEST", http.StatusBadRequest, message, nil)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError("NOT_FOUND", http.StatusNotFound, message, nil)
}

func NewConflictError(message string) *AppError {
	return NewAppError("CONFLICT", http.StatusConflict, message, nil)
}
