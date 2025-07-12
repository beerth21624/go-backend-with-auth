package domain_test

import (
	"testing"

	"venturex-backend/internal/app/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	t.Run("NewUUID", func(t *testing.T) {
		// Act
		id := domain.NewUUID()

		// Assert
		assert.NotEmpty(t, id)
		_, err := uuid.Parse(id.String())
		assert.NoError(t, err)
		assert.False(t, id.IsEmpty())
		assert.Equal(t, id.String(), id.Value())
	})

	t.Run("NewUUIDFromString", func(t *testing.T) {
		validUUID := uuid.New().String()
		testCases := []struct {
			name      string
			value     string
			want      domain.UUID
			expectErr error
		}{
			{"failure: empty string", "", "", domain.ErrEmptyUUID},
			{"failure: whitespace string", "   ", "", domain.ErrEmptyUUID},
			{"failure: invalid format", "not-a-uuid", "", domain.ErrInvalidUUID},
			{"success", validUUID, domain.UUID(validUUID), nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				got, err := domain.NewUUIDFromString(tc.value)

				// Assert
				if tc.expectErr != nil {
					assert.ErrorIs(t, err, tc.expectErr)
					assert.Empty(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.want, got)
					parsedWant, _ := uuid.Parse(tc.value)
					assert.Equal(t, parsedWant, got.UUID())
				}
			})
		}
	})
}

func TestUserID(t *testing.T) {
	t.Run("NewUserID", func(t *testing.T) {
		// Act
		id := domain.NewUserID()

		// Assert
		assert.NotEmpty(t, id)
		_, err := uuid.Parse(id.String())
		assert.NoError(t, err)
		assert.False(t, id.IsEmpty())
		assert.Equal(t, id.String(), id.Value())
	})

	t.Run("NewUserIDFromString", func(t *testing.T) {
		validUUID := uuid.New().String()
		testCases := []struct {
			name      string
			value     string
			want      domain.UserID
			expectErr error
		}{
			{"failure: invalid format", "not-a-uuid", "", domain.ErrInvalidUUID},
			{"success", validUUID, domain.UserID(validUUID), nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				got, err := domain.NewUserIDFromString(tc.value)

				// Assert
				if tc.expectErr != nil {
					assert.ErrorIs(t, err, tc.expectErr)
					assert.Empty(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.want, got)
				}
			})
		}
	})
}

func TestSessionID(t *testing.T) {
	t.Run("NewSessionID", func(t *testing.T) {
		// Act
		id := domain.NewSessionID()

		// Assert
		assert.NotEmpty(t, id)
		_, err := uuid.Parse(id.String())
		assert.NoError(t, err)
		assert.False(t, id.IsEmpty())
		assert.Equal(t, id.String(), id.Value())
	})

	t.Run("NewSessionIDFromString", func(t *testing.T) {
		validUUID := uuid.New().String()
		testCases := []struct {
			name      string
			value     string
			want      domain.SessionID
			expectErr error
		}{
			{"failure: invalid format", "not-a-uuid", "", domain.ErrInvalidUUID},
			{"success", validUUID, domain.SessionID(validUUID), nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Act
				got, err := domain.NewSessionIDFromString(tc.value)

				// Assert
				if tc.expectErr != nil {
					assert.ErrorIs(t, err, tc.expectErr)
					assert.Empty(t, got)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.want, got)
				}
			})
		}
	})
}
