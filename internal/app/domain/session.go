package domain

import (
	"time"
)

type Session struct {
	id                SessionID
	userID            UserID
	accessToken       JWT
	refreshTokenValue RefreshTokenValue
	deviceFingerprint DeviceFingerprint
	ipAddress         IPAddress
	userAgent         NonEmptyString
	isActive          bool
	expiresAt         Timestamp
	refreshExpiresAt  Timestamp
	createdAt         CreatedAt
	updatedAt         UpdatedAt
	lastActivity      Timestamp
}

func NewSession(
	userID UserID,
	accessToken string,
	deviceFingerprint, ipAddress, userAgent string,
	accessExpiresAt, refreshExpiresAt time.Time,
) (*Session, error) {

	accessTokenVO, err := NewJWT(accessToken)
	if err != nil {
		return nil, err
	}

	refreshTokenValue, err := GenerateRefreshTokenValue()
	if err != nil {
		return nil, err
	}

	deviceFingerprintVO, err := NewDeviceFingerprint(deviceFingerprint)
	if err != nil {
		return nil, err
	}

	ipAddressVO, err := NewIPAddress(ipAddress)
	if err != nil {
		return nil, err
	}

	userAgentVO, err := NewNonEmptyString(userAgent)
	if err != nil {
		return nil, err
	}

	expiresAtVO, err := NewTimestamp(accessExpiresAt)
	if err != nil {
		return nil, err
	}

	refreshExpiresAtVO, err := NewTimestamp(refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	createdAt, err := NewCreatedAt(now)
	if err != nil {
		return nil, err
	}

	updatedAt, err := NewUpdatedAt(now)
	if err != nil {
		return nil, err
	}

	lastActivityVO, err := NewTimestamp(now)
	if err != nil {
		return nil, err
	}

	return &Session{
		id:                NewSessionID(),
		userID:            userID,
		accessToken:       accessTokenVO,
		refreshTokenValue: refreshTokenValue,
		deviceFingerprint: deviceFingerprintVO,
		ipAddress:         ipAddressVO,
		userAgent:         userAgentVO,
		isActive:          true,
		expiresAt:         expiresAtVO,
		refreshExpiresAt:  refreshExpiresAtVO,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
		lastActivity:      lastActivityVO,
	}, nil
}

func ReconstructSession(
	id, userID, accessToken, refreshTokenValue, deviceFingerprint, ipAddress, userAgent string,
	isActive bool,
	expiresAt, refreshExpiresAt, createdAt, updatedAt, lastActivity time.Time,
) (*Session, error) {
	idVO, err := NewSessionIDFromString(id)
	if err != nil {
		return nil, err
	}

	userIdVO, err := NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}

	accessTokenVO, err := NewJWT(accessToken)
	if err != nil {
		return nil, err
	}

	refreshTokenValueVO, err := NewRefreshTokenValue(refreshTokenValue)
	if err != nil {
		return nil, err
	}

	deviceFingerprintVO, err := NewDeviceFingerprint(deviceFingerprint)
	if err != nil {
		return nil, err
	}

	ipAddressVO, err := NewIPAddress(ipAddress)
	if err != nil {
		return nil, err
	}

	userAgentVO, err := NewNonEmptyString(userAgent)
	if err != nil {
		return nil, err
	}

	expiresAtVO, err := NewTimestamp(expiresAt)
	if err != nil {
		return nil, err
	}

	refreshExpiresAtVO, err := NewTimestamp(refreshExpiresAt)
	if err != nil {
		return nil, err
	}

	createdAtVO, err := NewCreatedAt(createdAt)
	if err != nil {
		return nil, err
	}

	updatedAtVO, err := NewUpdatedAt(updatedAt)
	if err != nil {
		return nil, err
	}

	lastActivityVO, err := NewTimestamp(lastActivity)
	if err != nil {
		return nil, err
	}

	return &Session{
		id:                idVO,
		userID:            userIdVO,
		accessToken:       accessTokenVO,
		refreshTokenValue: refreshTokenValueVO,
		deviceFingerprint: deviceFingerprintVO,
		ipAddress:         ipAddressVO,
		userAgent:         userAgentVO,
		isActive:          isActive,
		expiresAt:         expiresAtVO,
		refreshExpiresAt:  refreshExpiresAtVO,
		createdAt:         createdAtVO,
		updatedAt:         updatedAtVO,
		lastActivity:      lastActivityVO,
	}, nil
}

// Getters
func (s *Session) ID() SessionID {
	return s.id
}

func (s *Session) UserID() UserID {
	return s.userID
}

func (s *Session) AccessToken() JWT {
	return s.accessToken
}

func (s *Session) RefreshTokenValue() RefreshTokenValue {
	return s.refreshTokenValue
}

func (s *Session) DeviceFingerprint() DeviceFingerprint {
	return s.deviceFingerprint
}

func (s *Session) IPAddress() IPAddress {
	return s.ipAddress
}

func (s *Session) UserAgent() NonEmptyString {
	return s.userAgent
}

func (s *Session) DeviceInfo() string {
	return s.userAgent.String()
}

func (s *Session) IsActive() bool {
	return s.isActive
}

func (s *Session) ExpiresAt() Timestamp {
	return s.expiresAt
}

func (s *Session) RefreshExpiresAt() Timestamp {
	return s.refreshExpiresAt
}

func (s *Session) CreatedAt() CreatedAt {
	return s.createdAt
}

func (s *Session) UpdatedAt() UpdatedAt {
	return s.updatedAt
}

func (s *Session) LastActivity() Timestamp {
	return s.lastActivity
}

// Business methods
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt.Time())
}

func (s *Session) IsRefreshExpired() bool {
	return time.Now().After(s.refreshExpiresAt.Time())
}

func (s *Session) IsValid() bool {
	return s.isActive && !s.IsExpired()
}

func (s *Session) CanRefresh() bool {
	return s.isActive && !s.IsRefreshExpired()
}

func (s *Session) Deactivate() {
	s.isActive = false
	s.updatedAt = NewUpdatedAtNow()
}

func (s *Session) RefreshAccessToken(newAccessToken string, newExpiresAt time.Time) error {
	if !s.CanRefresh() {
		return ErrRefreshTokenExpired
	}

	accessTokenVO, err := NewJWT(newAccessToken)
	if err != nil {
		return err
	}

	expiresAtVO, err := NewTimestamp(newExpiresAt)
	if err != nil {
		return err
	}

	s.accessToken = accessTokenVO
	s.expiresAt = expiresAtVO
	s.updatedAt = NewUpdatedAtNow()
	return nil
}

func (s *Session) MatchesDevice(deviceFingerprint, ipAddress string) bool {
	return s.deviceFingerprint.String() == deviceFingerprint &&
		s.ipAddress.String() == ipAddress
}

type LoginAttempt struct {
	id            ID
	username      NonEmptyString
	ipAddress     IPAddress
	userAgent     NonEmptyString
	success       bool
	failureReason *string
	attemptedAt   Timestamp
}

func NewLoginAttempt(
	username, ipAddress, userAgent string,
	success bool,
	failureReason *string,
) (*LoginAttempt, error) {
	usernameVO, err := NewNonEmptyString(username)
	if err != nil {
		return nil, err
	}

	ipAddressVO, err := NewIPAddress(ipAddress)
	if err != nil {
		return nil, err
	}

	userAgentVO, err := NewNonEmptyString(userAgent)
	if err != nil {
		return nil, err
	}

	attemptedAt, err := NewTimestamp(time.Now())
	if err != nil {
		return nil, err
	}

	return &LoginAttempt{
		username:      usernameVO,
		ipAddress:     ipAddressVO,
		userAgent:     userAgentVO,
		success:       success,
		failureReason: failureReason,
		attemptedAt:   attemptedAt,
	}, nil
}

func ReconstructLoginAttempt(
	id int64,
	username, ipAddress, userAgent string,
	success bool,
	failureReason *string,
	attemptedAt time.Time,
) (*LoginAttempt, error) {
	idVO, err := NewID(id)
	if err != nil {
		return nil, err
	}

	usernameVO, err := NewNonEmptyString(username)
	if err != nil {
		return nil, err
	}

	ipAddressVO, err := NewIPAddress(ipAddress)
	if err != nil {
		return nil, err
	}

	userAgentVO, err := NewNonEmptyString(userAgent)
	if err != nil {
		return nil, err
	}

	attemptedAtVO, err := NewTimestamp(attemptedAt)
	if err != nil {
		return nil, err
	}

	return &LoginAttempt{
		id:            idVO,
		username:      usernameVO,
		ipAddress:     ipAddressVO,
		userAgent:     userAgentVO,
		success:       success,
		failureReason: failureReason,
		attemptedAt:   attemptedAtVO,
	}, nil
}

func (la *LoginAttempt) ID() ID {
	return la.id
}

func (la *LoginAttempt) Username() NonEmptyString {
	return la.username
}

func (la *LoginAttempt) IPAddress() IPAddress {
	return la.ipAddress
}

func (la *LoginAttempt) UserAgent() NonEmptyString {
	return la.userAgent
}

func (la *LoginAttempt) Success() bool {
	return la.success
}

func (la *LoginAttempt) FailureReason() *string {
	return la.failureReason
}

func (la *LoginAttempt) AttemptedAt() Timestamp {
	return la.attemptedAt
}

func (la *LoginAttempt) IsRecent(within time.Duration) bool {
	return time.Since(la.attemptedAt.Time()) <= within
}

func (la *LoginAttempt) IsSuspicious() bool {
	return !la.success && la.failureReason != nil
}
