package domain

import "errors"

type ErrorCategory string

const (
	ErrCatValidation ErrorCategory = "validation"
	ErrCatBusiness   ErrorCategory = "business"
	ErrCatAuth       ErrorCategory = "auth"
	ErrCatSystem     ErrorCategory = "system"
)

type DomainError struct {
	Category ErrorCategory
	Code     string
	Message  string
	Err      error
}

func (e *DomainError) Error() string { return e.Message }
func (e *DomainError) Unwrap() error { return e.Err }

func (e *DomainError) Wrap(cause error) *DomainError {
	return &DomainError{
		Category: e.Category,
		Code:     e.Code,
		Message:  e.Message,
		Err:      cause,
	}
}

func (e *DomainError) Is(target error) bool {
	var t *DomainError
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

func DefineError(cat ErrorCategory, code, msg string) *DomainError {
	return &DomainError{
		Category: cat,
		Code:     code,
		Message:  msg,
	}
}

var (
	ErrUserNotFound          = DefineError(ErrCatBusiness, "USER_NOT_FOUND", "user not found")
	ErrInvalidCredentials    = DefineError(ErrCatAuth, "INVALID_CREDENTIALS", "invalid username or password")
	ErrAccountLocked         = DefineError(ErrCatBusiness, "ACCOUNT_LOCKED", "account is locked due to too many failed attempts")
	ErrTooManyFailedAttempts = DefineError(ErrCatBusiness, "TOO_MANY_FAILED_ATTEMPTS", "too many failed login attempts")
	ErrTokenExpired          = DefineError(ErrCatAuth, "TOKEN_EXPIRED", "token has expired")
	ErrTokenInvalid          = DefineError(ErrCatAuth, "TOKEN_INVALID", "token is invalid")
	ErrSessionNotFound       = DefineError(ErrCatAuth, "SESSION_NOT_FOUND", "session not found")
	ErrInvalidSession        = DefineError(ErrCatAuth, "INVALID_SESSION", "session is invalid")
	ErrRefreshTokenExpired   = DefineError(ErrCatAuth, "REFRESH_TOKEN_EXPIRED", "refresh token has expired")
)
