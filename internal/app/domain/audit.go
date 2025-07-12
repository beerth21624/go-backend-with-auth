package domain

import (
	"errors"
	"strings"
)

const (
	AuditUserSystem = AuditUser("System")
)

var (
	ErrInvalidAuditUser = errors.New("invalid audit user")
)

type AuditUser string

func NewAuditUser(value string) (AuditUser, error) {
	if value == "" || strings.TrimSpace(value) != value {
		var zero AuditUser
		return zero, ErrInvalidAuditUser
	}
	return AuditUser(value), nil
}

func (a AuditUser) String() string {
	return string(a)
}

type Audit struct {
	user AuditUser
	date Timestamp
}

func (a Audit) User() AuditUser { return a.user }
func (a Audit) Date() Timestamp { return a.date }

func NewAudit(user AuditUser, date Timestamp) Audit {
	return Audit{user: user, date: date}
}

func (a Audit) WithDate(date Timestamp) Audit {
	return Audit{
		user: a.user,
		date: date,
	}
}

func (a Audit) WithUser(user AuditUser) Audit {
	return Audit{
		user: user,
		date: a.date,
	}
}
