package api

import (
	"errors"

	"venturex-backend/internal/app/domain"
	"venturex-backend/internal/pkg/validator"
)

func MapDomainError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	var validationErr *validator.ValidationError
	if errors.As(err, &validationErr) {
		return NewValidationError(validationErr).WithCause(err)
	}

	var de *domain.DomainError
	if errors.As(err, &de) {
		switch de.Category {
		case domain.ErrCatValidation:
			return NewBadRequestError(de.Message).WithCause(err)
		case domain.ErrCatAuth:
			return NewUnauthorizedError(de.Message).WithCause(err)
		case domain.ErrCatBusiness:
			return NewConflictError(de.Message).WithCause(err)
		case domain.ErrCatSystem:
			return NewInternalError(err)
		}
	}

	return NewInternalError(err)
}

func (e *AppError) WithCause(cause error) *AppError {
	e.Err = cause
	return e
}
