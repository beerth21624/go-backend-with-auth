package domain_test

import (
	"testing"

	"venturex-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewNonEmptyString(t *testing.T) {
	// Arrange
	testCases := []struct {
		name      string
		value     string
		want      domain.NonEmptyString
		expectErr error
	}{
		{
			name:      "failure: empty string",
			value:     "",
			want:      "",
			expectErr: domain.ErrEmptyString,
		},
		{
			name:      "failure: string with only whitespace",
			value:     "   ",
			want:      "",
			expectErr: domain.ErrEmptyString,
		},
		{
			name:      "success: valid non-empty string",
			value:     "hello world",
			want:      "hello world",
			expectErr: nil,
		},
		{
			name:      "success: string with leading/trailing whitespace",
			value:     "  padded string  ",
			want:      "  padded string  ",
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := domain.NewNonEmptyString(tc.value)

			// Assert
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.String())
				assert.Equal(t, tc.value, got.Value())
			}
		})
	}
}
