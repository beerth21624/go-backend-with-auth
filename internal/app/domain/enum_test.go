package domain_test

import (
	"errors"
	"testing"

	"venturex-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewStatus(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		want      domain.Status
		expectErr error
	}{
		{"failure: invalid status", "unknown", "", domain.ErrInvalidStatus},
		{"success: active", "active", domain.StatusActive, nil},
		{"success: inactive", "inactive", domain.StatusInactive, nil},
		{"success: pending", "pending", domain.StatusPending, nil},
		{"success: deleted", "deleted", domain.StatusDeleted, nil},
		{"success: case-insensitivity", "  aCTiVe  ", domain.StatusActive, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewStatus(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, string(tc.want), got.String())
				assert.Equal(t, string(tc.want), got.Value())
			}
		})
	}
}

func TestStatusMethods(t *testing.T) {
	assert.True(t, domain.StatusActive.IsActive())
	assert.False(t, domain.StatusInactive.IsActive())

	assert.True(t, domain.StatusInactive.IsInactive())
	assert.False(t, domain.StatusActive.IsInactive())

	assert.True(t, domain.StatusPending.IsPending())
	assert.False(t, domain.StatusActive.IsPending())

	assert.True(t, domain.StatusDeleted.IsDeleted())
	assert.False(t, domain.StatusActive.IsDeleted())
}

func TestNewUserRole(t *testing.T) {
	errInvalidUserRole := errors.New("invalid user role")

	testCases := []struct {
		name      string
		value     string
		want      domain.UserRole
		expectErr error
	}{
		{"failure: invalid role", "superadmin", "", errInvalidUserRole},
		{"success: admin", "admin", domain.UserRoleAdmin, nil},
		{"success: user", "user", domain.UserRoleUser, nil},
		{"success: guest", "guest", domain.UserRoleGuest, nil},
		{"success: case-insensitivity", "  aDmIn  ", domain.UserRoleAdmin, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewUserRole(tc.value)
			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, string(tc.want), got.String())
				assert.Equal(t, string(tc.want), got.Value())
			}
		})
	}
}

func TestUserRoleMethods(t *testing.T) {
	assert.True(t, domain.UserRoleAdmin.IsAdmin())
	assert.False(t, domain.UserRoleUser.IsAdmin())

	assert.True(t, domain.UserRoleUser.IsUser())
	assert.False(t, domain.UserRoleAdmin.IsUser())

	assert.True(t, domain.UserRoleGuest.IsGuest())
	assert.False(t, domain.UserRoleAdmin.IsGuest())
}
