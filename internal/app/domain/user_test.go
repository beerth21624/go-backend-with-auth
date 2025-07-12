package domain_test

import (
	"testing"
	"time"

	"beerdosan-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	// Arrange
	testCases := []struct {
		name          string
		username      string
		email         string
		firstName     string
		lastName      string
		plainPassword string
		expectErr     bool
	}{
		{"success", "testuser", "test@example.com", "Test", "User", "password123", false},
		{"failure: empty username", "", "test@example.com", "Test", "User", "password123", true},
		{"failure: short password", "testuser", "test@example.com", "Test", "User", "123", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			user, err := domain.NewUser(tc.username, tc.email, tc.firstName, tc.lastName, tc.plainPassword)

			// Assert
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, domain.NonEmptyString(tc.username), user.Username())
				assert.Equal(t, domain.NonEmptyString(tc.email), user.Email())
				assert.Equal(t, tc.firstName+" "+tc.lastName, user.FullName())
				assert.True(t, user.VerifyPassword(tc.plainPassword))
				assert.True(t, user.IsActive())
				assert.Equal(t, domain.UserRoleUser, user.Role())
			}
		})
	}
}

func TestReconstructUser(t *testing.T) {
	// Arrange
	userID := domain.NewUserID().String()
	hashedPassword, err := domain.NewHashedPassword("password123")
	require.NoError(t, err)
	now := time.Now()

	// Act
	user, err := domain.ReconstructUser(
		userID,
		"re_user",
		"re@example.com",
		"Re",
		"Construct",
		hashedPassword.String(),
		"admin",
		"inactive",
		now,
		now,
	)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, domain.UserID(userID), user.ID())
	assert.Equal(t, domain.UserRoleAdmin, user.Role())
	assert.Equal(t, domain.StatusInactive, user.Status())
	assert.False(t, user.IsActive())
}

func TestUser_BusinessLogic(t *testing.T) {
	// Arrange
	user, err := domain.NewUser("testuser", "test@example.com", "Test", "User", "password123")
	require.NoError(t, err)

	t.Run("Activate/Deactivate", func(t *testing.T) {
		// Act & Assert
		err := user.Deactivate()
		require.NoError(t, err)
		assert.False(t, user.IsActive())
		assert.Equal(t, domain.StatusInactive, user.Status())
		originalUpdatedAt := user.UpdatedAt()

		time.Sleep(1 * time.Millisecond) // Ensure UpdatedAt changes

		err = user.Activate()
		require.NoError(t, err)
		assert.True(t, user.IsActive())
		assert.Equal(t, domain.StatusActive, user.Status())
		assert.NotEqual(t, originalUpdatedAt, user.UpdatedAt())
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		// Act
		originalUpdatedAt := user.UpdatedAt()
		time.Sleep(1 * time.Millisecond)
		err := user.UpdateProfile("NewFirst", "NewLast")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "NewFirst", user.FirstName().String())
		assert.Equal(t, "NewLast", user.LastName().String())
		assert.Equal(t, "NewFirst NewLast", user.FullName())
		assert.NotEqual(t, originalUpdatedAt, user.UpdatedAt())
	})

	t.Run("ChangePassword", func(t *testing.T) {
		// Act
		originalUpdatedAt := user.UpdatedAt()
		time.Sleep(1 * time.Millisecond)
		err := user.ChangePassword("newPassword456")

		// Assert
		require.NoError(t, err)
		assert.True(t, user.VerifyPassword("newPassword456"))
		assert.False(t, user.VerifyPassword("password123"))
		assert.NotEqual(t, originalUpdatedAt, user.UpdatedAt())

		// Test failure
		err = user.ChangePassword("short")
		assert.Error(t, err)
	})
}
