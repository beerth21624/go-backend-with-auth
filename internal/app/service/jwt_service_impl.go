package service

import (
	"crypto/sha256"
	"encoding/binary"
	"strconv"
	"beerdosan-backend/internal/app/domain"
	"beerdosan-backend/internal/pkg/jwt"
)

type jwtServiceImpl struct {
	jwtService jwt.JWTService
}

func NewJWTService(jwtService jwt.JWTService) JWTService {
	return &jwtServiceImpl{
		jwtService: jwtService,
	}
}

func uuidToInt64(uuid string) int64 {
	hash := sha256.Sum256([]byte(uuid))
	return int64(binary.BigEndian.Uint64(hash[:8]))
}

func (s *jwtServiceImpl) GenerateAccessToken(userID domain.UserID, sessionID domain.SessionID, role domain.UserRole) (domain.JWT, error) {
	userIDInt64 := uuidToInt64(userID.String())
	sessionIDInt64 := uuidToInt64(sessionID.String())

	token, _, err := s.jwtService.GenerateAccessToken(userIDInt64, userID.String(), "", sessionIDInt64, sessionID.String())
	if err != nil {
		return "", err
	}

	return domain.JWT(token), nil
}

func (s *jwtServiceImpl) GenerateRefreshToken(userID domain.UserID, sessionID domain.SessionID) (domain.JWT, error) {
	userIDInt64 := uuidToInt64(userID.String())
	sessionIDInt64 := uuidToInt64(sessionID.String())

	token, _, err := s.jwtService.GenerateRefreshToken(userIDInt64, userID.String(), sessionID.String(), sessionIDInt64)
	if err != nil {
		return "", err
	}

	return domain.JWT(token), nil
}

func (s *jwtServiceImpl) ValidateToken(token domain.JWT) (*domain.TokenClaims, error) {
	claims, err := s.jwtService.ValidateToken(token.String())
	if err != nil {
		return nil, err
	}

	userID := claims.Username
	sessionID := claims.Fingerprint

	if sessionID == "" {
		sessionID = claims.Email
	}

	if userID == "" {
		userID = strconv.FormatInt(claims.UserID, 10)
	}
	if sessionID == "" {
		sessionID = strconv.FormatInt(claims.SessionID, 10)
	}

	return &domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Role:      claims.TokenType,
		TokenType: claims.TokenType,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

func (s *jwtServiceImpl) RefreshAccessToken(refreshToken domain.JWT) (domain.JWT, error) {
	token, _, err := s.jwtService.RefreshAccessToken(refreshToken.String())
	if err != nil {
		return "", err
	}

	return domain.JWT(token), nil
}

func (s *jwtServiceImpl) RevokeToken(token domain.JWT) error {
	// TODO: Implement JWT revocation
	// For now, JWT revocation is not implemented in this service
	// This would typically involve adding the token to a blacklist
	// or updating the session status
	return nil
}
