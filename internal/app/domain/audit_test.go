package domain_test

import (
	"testing"
	"time"

	"venturex-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditUser(t *testing.T) {
	// Arrange
	testCases := []struct {
		name      string
		value     string
		want      domain.AuditUser
		expectErr error
	}{
		{
			name:      "failure: empty value",
			value:     "",
			want:      "",
			expectErr: domain.ErrInvalidAuditUser,
		},
		{
			name:      "failure: value with only whitespace",
			value:     "   ",
			want:      "",
			expectErr: domain.ErrInvalidAuditUser,
		},
		{
			name:      "failure: value with leading/trailing whitespace",
			value:     " user ",
			want:      "",
			expectErr: domain.ErrInvalidAuditUser,
		},
		{
			name:      "success: valid value",
			value:     "test_user",
			want:      domain.AuditUser("test_user"),
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := domain.NewAuditUser(tc.value)

			// Assert
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.String())
			}
		})
	}
}

func TestAudit(t *testing.T) {
	// Arrange
	baseUser, err := domain.NewAuditUser("base_user")
	require.NoError(t, err)
	baseTime := time.Now()
	baseTs, err := domain.NewTimestamp(baseTime)
	require.NoError(t, err)

	type arrangeFunc func(t *testing.T) (initial domain.Audit, result domain.Audit)
	type assertFunc func(t *testing.T, initial domain.Audit, result domain.Audit)

	testCases := []struct {
		name    string
		arrange arrangeFunc
		assert  assertFunc
	}{
		{
			name: "success: create new audit",
			arrange: func(t *testing.T) (domain.Audit, domain.Audit) {
				audit := domain.NewAudit(baseUser, baseTs)
				return audit, audit
			},
			assert: func(t *testing.T, initial domain.Audit, result domain.Audit) {
				assert.Equal(t, baseUser, result.User())
				assert.Equal(t, baseTs, result.Date())
			},
		},
		{
			name: "success: WithDate should not mutate original",
			arrange: func(t *testing.T) (domain.Audit, domain.Audit) {
				initialAudit := domain.NewAudit(baseUser, baseTs)
				newTime := baseTime.Add(1 * time.Hour)
				newTs, err := domain.NewTimestamp(newTime)
				require.NoError(t, err)

				resultAudit := initialAudit.WithDate(newTs)
				return initialAudit, resultAudit
			},
			assert: func(t *testing.T, initial domain.Audit, result domain.Audit) {
				newTime := baseTime.Add(1 * time.Hour)
				newTs, err := domain.NewTimestamp(newTime)
				require.NoError(t, err)

				// Assert result
				assert.Equal(t, baseUser, result.User())
				assert.Equal(t, newTs, result.Date())

				// Assert original is not mutated
				assert.Equal(t, baseTs, initial.Date())
			},
		},
		{
			name: "success: WithUser should not mutate original",
			arrange: func(t *testing.T) (domain.Audit, domain.Audit) {
				initialAudit := domain.NewAudit(baseUser, baseTs)
				newUser, err := domain.NewAuditUser("new_user")
				require.NoError(t, err)

				resultAudit := initialAudit.WithUser(newUser)
				return initialAudit, resultAudit
			},
			assert: func(t *testing.T, initial domain.Audit, result domain.Audit) {
				newUser, err := domain.NewAuditUser("new_user")
				require.NoError(t, err)

				// Assert result
				assert.Equal(t, newUser, result.User())
				assert.Equal(t, baseTs, result.Date())

				// Assert original is not mutated
				assert.Equal(t, baseUser, initial.User())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			initial, result := tc.arrange(t)

			// Assert
			tc.assert(t, initial, result)
		})
	}
}
