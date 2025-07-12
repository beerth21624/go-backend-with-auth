package domain_test

import (
	"testing"
	"time"

	"beerdosan-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashedPassword(t *testing.T) {
	t.Run("NewHashedPassword", func(t *testing.T) {
		// Valid password
		hp, err := domain.NewHashedPassword("password123")
		require.NoError(t, err)
		assert.NotEmpty(t, hp)
		assert.True(t, hp.VerifyPassword("password123"))
		assert.False(t, hp.VerifyPassword("wrongpassword"))
		assert.NotEmpty(t, hp.String())

		// Invalid password (too short)
		_, err = domain.NewHashedPassword("short")
		assert.ErrorIs(t, err, domain.ErrInvalidPassword)
	})

	t.Run("NewHashedPasswordFromHash", func(t *testing.T) {
		hp, err := domain.NewHashedPasswordFromHash("$2a$10$.....................................................")
		require.NoError(t, err)
		assert.NotEmpty(t, hp)

		_, err = domain.NewHashedPasswordFromHash("")
		assert.ErrorIs(t, err, domain.ErrInvalidPassword)
	})
}

func TestJWT(t *testing.T) {
	validJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	testCases := []struct {
		name      string
		value     string
		want      domain.JWT
		expectErr error
	}{
		{"failure: empty", "", "", domain.ErrInvalidJWT},
		{"failure: whitespace", "  ", "", domain.ErrInvalidJWT},
		{"failure: invalid format", "not-a-valid-jwt-string", "", domain.ErrInvalidJWT},
		{"success", validJWT, domain.JWT(validJWT), nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewJWT(tc.value)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, tc.value, got.String())
				assert.False(t, got.IsEmpty())
			}
		})
	}
}

func TestIPAddress(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		isPrivate bool
		expectErr bool
	}{
		{"failure: empty", "", false, true},
		{"failure: invalid", "not an ip", false, true},
		{"success: public ipv4", "8.8.8.8", false, false},
		{"success: private ipv4", "192.168.1.1", true, false},
		{"success: localhost", "127.0.0.1", true, false},
		{"success: public ipv6", "2001:4860:4860::8888", false, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewIPAddress(tc.value)
			if tc.expectErr {
				assert.ErrorIs(t, err, domain.ErrInvalidIPAddress)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, domain.IPAddress(tc.value), got)
				assert.Equal(t, tc.isPrivate, got.IsPrivate())
			}
		})
	}
}

func TestTokenType(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		want      domain.TokenType
		expectErr bool
	}{
		{"failure: invalid", "bearer", "", true},
		{"success: access", "access", domain.TokenTypeAccess, false},
		{"success: refresh", "refresh", domain.TokenTypeRefresh, false},
		{"success: case-insensitive", "  aCCeSS  ", domain.TokenTypeAccess, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewTokenType(tc.value)
			if tc.expectErr {
				assert.ErrorIs(t, err, domain.ErrInvalidTokenType)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
				assert.Equal(t, got.IsAccess(), tc.want == domain.TokenTypeAccess)
				assert.Equal(t, got.IsRefresh(), tc.want == domain.TokenTypeRefresh)
			}
		})
	}
}

func TestDeviceFingerprint(t *testing.T) {
	t.Run("NewDeviceFingerprint", func(t *testing.T) {
		_, err := domain.NewDeviceFingerprint("short")
		assert.ErrorIs(t, err, domain.ErrInvalidFingerprint)

		longString := string(make([]byte, 501))
		_, err = domain.NewDeviceFingerprint(longString)
		assert.ErrorIs(t, err, domain.ErrInvalidFingerprint)

		fp, err := domain.NewDeviceFingerprint("valid-fingerprint-of-appropriate-length")
		require.NoError(t, err)
		assert.Equal(t, "valid-fingerprint-of-appropriate-length", fp.String())
	})

	t.Run("GenerateDeviceFingerprint", func(t *testing.T) {
		_, err := domain.GenerateDeviceFingerprint("", "1.1.1.1")
		assert.ErrorIs(t, err, domain.ErrInvalidFingerprint)

		fp, err := domain.GenerateDeviceFingerprint("my-user-agent", "127.0.0.1")
		require.NoError(t, err)
		assert.NotEmpty(t, fp)
	})
}

func TestRefreshTokenValue(t *testing.T) {
	t.Run("NewRefreshTokenValue", func(t *testing.T) {
		_, err := domain.NewRefreshTokenValue("too-short")
		assert.ErrorIs(t, err, domain.ErrInvalidJWT)

		val := "a-sufficiently-long-refresh-token-value-that-is-at-least-32-chars"
		rtv, err := domain.NewRefreshTokenValue(val)
		require.NoError(t, err)
		assert.Equal(t, val, rtv.String())
		assert.False(t, rtv.IsEmpty())
	})

	t.Run("GenerateRefreshTokenValue", func(t *testing.T) {
		rtv, err := domain.GenerateRefreshTokenValue()
		require.NoError(t, err)
		assert.NotEmpty(t, rtv)
		assert.False(t, rtv.IsEmpty())
	})
}

func TestTokenClaims(t *testing.T) {
	claims := domain.TokenClaims{
		TokenType: string(domain.TokenTypeAccess),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	assert.False(t, claims.IsExpired())
	assert.True(t, claims.IsAccessToken())
	assert.False(t, claims.IsRefreshToken())

	claims.ExpiresAt = time.Now().Add(-1 * time.Hour)
	assert.True(t, claims.IsExpired())

	claims.TokenType = string(domain.TokenTypeRefresh)
	assert.False(t, claims.IsAccessToken())
	assert.True(t, claims.IsRefreshToken())
}
