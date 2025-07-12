package domain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyString    = errors.New("string cannot be empty")
	ErrStringTooLong  = errors.New("string too long")
	ErrStringTooShort = errors.New("string too short")
)

type NonEmptyString string

func NewNonEmptyString(s string) (NonEmptyString, error) {
	if strings.TrimSpace(s) == "" {
		return "", ErrEmptyString
	}
	return NonEmptyString(s), nil
}

func (s NonEmptyString) String() string {
	return string(s)
}

func (s NonEmptyString) Value() string {
	return string(s)
}
