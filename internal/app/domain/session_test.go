package domain_test

import (
	"testing"
	"time"

	"beerdosan-backend/internal/app/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	validRefreshToken = "valid-refresh-token-that-is-long-enough-for-validation-to-pass"
)

func TestSession(t *testing.T) {
	userID := domain.NewUserID()
	fp, err := domain.GenerateDeviceFingerprint("ua", "127.0.0.1")
	require.NoError(t, err)

	t.Run("NewSession", func(t *testing.T) {
		// Arrange
		accessExpiresAt := time.Now().Add(15 * time.Minute)
		refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

		// Act
		session, err := domain.NewSession(
			userID,
			validAccessToken,
			string(fp),
			"127.0.0.1",
			"test-user-agent",
			accessExpiresAt,
			refreshExpiresAt,
		)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserID())
		assert.True(t, session.IsActive())
		assert.False(t, session.IsExpired())
		assert.False(t, session.IsRefreshExpired())
		assert.True(t, session.IsValid())
		assert.True(t, session.CanRefresh())
		assert.Equal(t, "test-user-agent", session.DeviceInfo())
	})

	t.Run("ReconstructSession", func(t *testing.T) {
		// Arrange
		sessionID := domain.NewSessionID()
		now := time.Now()

		// Act
		session, err := domain.ReconstructSession(
			sessionID.String(),
			userID.String(),
			validAccessToken,
			validRefreshToken,
			string(fp),
			"192.168.1.1",
			"reconstructed-ua",
			true,
			now.Add(1*time.Hour),
			now.Add(2*time.Hour),
			now,
			now,
			now,
		)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, sessionID, session.ID())
		assert.Equal(t, userID, session.UserID())
		assert.True(t, session.IsActive())
	})

	t.Run("Session business methods", func(t *testing.T) {
		// Arrange
		activeSession, _ := domain.NewSession(userID, validAccessToken, string(fp), "1.1.1.1", "ua", time.Now().Add(1*time.Hour), time.Now().Add(2*time.Hour))
		expiredSession, _ := domain.NewSession(userID, validAccessToken, string(fp), "1.1.1.1", "ua", time.Now().Add(-1*time.Hour), time.Now().Add(2*time.Hour))
		expiredRefreshSession, _ := domain.NewSession(userID, validAccessToken, string(fp), "1.1.1.1", "ua", time.Now().Add(-2*time.Hour), time.Now().Add(-1*time.Hour))

		// Deactivate
		activeSession.Deactivate()
		assert.False(t, activeSession.IsActive())

		// IsExpired & IsValid
		assert.True(t, expiredSession.IsExpired())
		assert.False(t, expiredSession.IsValid())

		// IsRefreshExpired & CanRefresh
		assert.True(t, expiredRefreshSession.IsRefreshExpired())
		assert.False(t, expiredRefreshSession.CanRefresh())

		// RefreshAccessToken
		newAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJzYW1wbGV0ZXN0In0.j-mG9R_p_k2uqK2aT8cbH-h_L6PjXwG2vJ-2aR2bY_c"
		err := expiredSession.RefreshAccessToken(newAccessToken, time.Now().Add(1*time.Hour))
		require.NoError(t, err)
		assert.Equal(t, domain.JWT(newAccessToken), expiredSession.AccessToken())
		assert.False(t, expiredSession.IsExpired())

		// MatchesDevice
		assert.True(t, expiredSession.MatchesDevice(string(fp), "1.1.1.1"))
		assert.False(t, expiredSession.MatchesDevice("other-fp", "1.1.1.1"))
	})
}

func TestLoginAttempt(t *testing.T) {
	username, _ := domain.NewNonEmptyString("testuser")
	ip, _ := domain.NewIPAddress("127.0.0.1")
	ua, _ := domain.NewNonEmptyString("test-agent")

	t.Run("NewLoginAttempt", func(t *testing.T) {
		// Act
		failureReason := "wrong password"
		suspiciousAttempt, err := domain.NewLoginAttempt(
			username.String(),
			ip.String(),
			ua.String(),
			false,
			&failureReason,
		)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, suspiciousAttempt)
		assert.Equal(t, username, suspiciousAttempt.Username())
		assert.False(t, suspiciousAttempt.Success())
		assert.Equal(t, &failureReason, suspiciousAttempt.FailureReason())
		assert.True(t, suspiciousAttempt.IsRecent(1*time.Minute))
		assert.True(t, suspiciousAttempt.IsSuspicious()) // Should be suspicious

		// Test non-suspicious attempt (success)
		successfulAttempt, err := domain.NewLoginAttempt(
			username.String(),
			ip.String(),
			ua.String(),
			true,
			nil,
		)
		require.NoError(t, err)
		assert.False(t, successfulAttempt.IsSuspicious())

		// Test non-suspicious attempt (failure but no reason)
		failedAttemptNoReason, err := domain.NewLoginAttempt(
			username.String(),
			ip.String(),
			ua.String(),
			false,
			nil,
		)
		require.NoError(t, err)
		assert.False(t, failedAttemptNoReason.IsSuspicious())
	})

	t.Run("ReconstructLoginAttempt", func(t *testing.T) {
		// Arrange
		now := time.Now()
		reason := "reconstructed reason"

		// Act
		attempt, err := domain.ReconstructLoginAttempt(
			1,
			"user",
			"1.1.1.1",
			"ua",
			true,
			&reason,
			now,
		)

		// Assert
		require.NoError(t, err)
		id, _ := domain.NewID(1)
		assert.Equal(t, id, attempt.ID())
		assert.True(t, attempt.Success())
		assert.Equal(t, &reason, attempt.FailureReason())
	})
}
