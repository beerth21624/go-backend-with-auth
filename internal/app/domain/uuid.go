package domain

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidUUID = errors.New("invalid UUID format")
	ErrEmptyUUID   = errors.New("UUID cannot be empty")
)

type UUID string

func NewUUID() UUID {
	return UUID(uuid.New().String())
}

func NewUUIDFromString(s string) (UUID, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ErrEmptyUUID
	}

	if _, err := uuid.Parse(s); err != nil {
		return "", ErrInvalidUUID
	}

	return UUID(s), nil
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) Value() string {
	return string(u)
}

func (u UUID) IsEmpty() bool {
	return string(u) == ""
}

func (u UUID) UUID() uuid.UUID {
	parsed, _ := uuid.Parse(string(u))
	return parsed
}

// UserID represents a user identifier
type UserID UUID

func NewUserID() UserID {
	return UserID(NewUUID())
}

func NewUserIDFromString(s string) (UserID, error) {
	uuid, err := NewUUIDFromString(s)
	if err != nil {
		return "", err
	}
	return UserID(uuid), nil
}

func (u UserID) String() string {
	return string(u)
}

func (u UserID) Value() string {
	return string(u)
}

func (u UserID) IsEmpty() bool {
	return string(u) == ""
}

type SessionID UUID

func NewSessionID() SessionID {
	return SessionID(NewUUID())
}

func NewSessionIDFromString(s string) (SessionID, error) {
	uuid, err := NewUUIDFromString(s)
	if err != nil {
		return "", err
	}
	return SessionID(uuid), nil
}

func (s SessionID) String() string {
	return string(s)
}

func (s SessionID) Value() string {
	return string(s)
}

func (s SessionID) IsEmpty() bool {
	return string(s) == ""
}
