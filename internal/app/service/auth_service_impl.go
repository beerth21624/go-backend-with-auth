package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"beerdosan-backend/internal/app/domain"
)

type AuthClaims struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	SessionID   int64  `json:"session_id"`
	Role        string `json:"role"`
	UserUUID    string `json:"user_uuid"`
	SessionUUID string `json:"session_uuid"`
}

func (s *AuthServiceImpl) ValidateCredentials(ctx context.Context, username, password string) (*domain.User, error) {
	log.Printf("[DEBUG] ValidateCredentials called: username=%s", username)
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.VerifyPassword(password) {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthServiceImpl) CreateSession(ctx context.Context, userID domain.UserID, deviceInfo, ipAddress string) (*domain.Session, error) {
	log.Printf("[DEBUG] CreateSession called: userID=%s deviceInfo=%s ipAddress=%s", userID, deviceInfo, ipAddress)

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for session creation: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	deviceFingerprint, err := domain.GenerateDeviceFingerprint(deviceInfo, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to generate device fingerprint: %w", err)
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	sessionID := domain.NewSessionID()

	accessToken, err := s.jwtService.GenerateAccessToken(userID, sessionID, user.Role())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	session, err := domain.ReconstructSession(
		sessionID.String(),
		userID.String(),
		accessToken.String(),
		refreshToken.String(),
		deviceFingerprint.String(),
		ipAddress,
		deviceInfo,
		true,
		expiresAt,
		refreshExpiresAt,
		time.Now(),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	createdSession, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return createdSession, nil
}

func (s *AuthServiceImpl) ValidateSession(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	log.Printf("[DEBUG] ValidateSession called: sessionID=%s", sessionID)
	session, err := s.sessionRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, domain.ErrSessionNotFound
	}

	if session == nil || !session.IsValid() {
		return nil, domain.ErrSessionNotFound
	}

	return session, nil
}

func (s *AuthServiceImpl) InvalidateSession(ctx context.Context, sessionID domain.SessionID) error {
	log.Printf("[DEBUG] InvalidateSession called: sessionID=%s", sessionID)
	return s.sessionRepo.InvalidateSession(ctx, sessionID)
}

func (s *AuthServiceImpl) InvalidateAllUserSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error {
	log.Printf("[DEBUG] InvalidateAllUserSessions called: userID=%s excludeSessionID=%s", userID, excludeSessionID)
	return s.sessionRepo.InvalidateAllUserSessions(ctx, userID, excludeSessionID)
}

func (s *AuthServiceImpl) UpdateSessionActivity(ctx context.Context, sessionID domain.SessionID) error {
	log.Printf("[DEBUG] UpdateSessionActivity called: sessionID=%s", sessionID)
	return s.sessionRepo.UpdateLastActivity(ctx, sessionID)
}

func (s *AuthServiceImpl) RecordLoginAttempt(ctx context.Context, username, ipAddress string, success bool, reason string) error {
	log.Printf("[DEBUG] RecordLoginAttempt called: username=%s ipAddress=%s success=%v reason=%s", username, ipAddress, success, reason)
	var failureReasonPtr *string
	if reason != "" {
		failureReasonPtr = &reason
	}

	attempt, err := domain.NewLoginAttempt(
		username,
		ipAddress,
		"unknown", // User agent is not available in this context
		success,
		failureReasonPtr,
	)
	if err != nil {
		return fmt.Errorf("failed to create login attempt: %w", err)
	}

	return s.loginAttemptRepo.Create(ctx, attempt)
}

func (s *AuthServiceImpl) CheckRateLimit(ctx context.Context, username, ipAddress string) error {
	log.Printf("[DEBUG] CheckRateLimit called: username=%s ipAddress=%s", username, ipAddress)
	since := time.Now().Add(-15 * time.Minute).Unix()

	failedCount, err := s.loginAttemptRepo.CountFailedAttemptsByUsernameAndIP(ctx, username, ipAddress, time.Unix(since, 0))
	if err != nil {
		return fmt.Errorf("failed to count failed attempts: %w", err)
	}

	if failedCount >= 5 {
		return domain.ErrAccountLocked
	}

	return nil
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (*AuthClaims, error) {
	jwtToken, err := domain.NewJWT(token)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	claims, err := s.jwtService.ValidateToken(jwtToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	if claims.IsExpired() {
		return nil, domain.ErrTokenExpired
	}

	userID := domain.UserID(claims.UserID)
	sessionID := domain.SessionID(claims.SessionID)

	userDomain, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	if userDomain == nil {
		return nil, domain.ErrTokenInvalid
	}

	sessionDomain, err := s.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, domain.ErrInvalidSession
	}

	if sessionDomain.UserID().String() != userID.String() {
		return nil, domain.ErrInvalidSession
	}

	userIDInt64 := uuidToInt64(userID.String())
	sessionIDInt64 := uuidToInt64(sessionID.String())

	return &AuthClaims{
		UserID:      userIDInt64,
		Username:    userDomain.Username().String(),
		Email:       userDomain.Email().String(),
		SessionID:   sessionIDInt64,
		Role:        claims.Role,
		UserUUID:    userID.String(),
		SessionUUID: sessionID.String(),
	}, nil
}
