package domain

import (
	"errors"
	"strings"
)

var (
	ErrInvalidStatus = errors.New("invalid status")
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
	StatusDeleted  Status = "deleted"
)

func NewStatus(s string) (Status, error) {
	status := Status(strings.ToLower(strings.TrimSpace(s)))
	switch status {
	case StatusActive, StatusInactive, StatusPending, StatusDeleted:
		return status, nil
	default:
		return "", ErrInvalidStatus
	}
}

func (s Status) String() string {
	return string(s)
}

func (s Status) Value() string {
	return string(s)
}

func (s Status) IsActive() bool {
	return s == StatusActive
}

func (s Status) IsInactive() bool {
	return s == StatusInactive
}

func (s Status) IsPending() bool {
	return s == StatusPending
}

func (s Status) IsDeleted() bool {
	return s == StatusDeleted
}

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
	UserRoleGuest UserRole = "guest"
)

func NewUserRole(s string) (UserRole, error) {
	role := UserRole(strings.ToLower(strings.TrimSpace(s)))
	switch role {
	case UserRoleAdmin, UserRoleUser, UserRoleGuest:
		return role, nil
	default:
		return "", errors.New("invalid user role")
	}
}

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) Value() string {
	return string(r)
}

func (r UserRole) IsAdmin() bool {
	return r == UserRoleAdmin
}

func (r UserRole) IsUser() bool {
	return r == UserRoleUser
}

func (r UserRole) IsGuest() bool {
	return r == UserRoleGuest
}
