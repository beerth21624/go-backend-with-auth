package service

import (
	"beerdosan-backend/internal/app/domain"
)

type JWTService interface {
	GenerateAccessToken(userID domain.UserID, sessionID domain.SessionID, role domain.UserRole) (domain.JWT, error)
	GenerateRefreshToken(userID domain.UserID, sessionID domain.SessionID) (domain.JWT, error)
	ValidateToken(token domain.JWT) (*domain.TokenClaims, error)
	RefreshAccessToken(refreshToken domain.JWT) (domain.JWT, error)
	RevokeToken(token domain.JWT) error
}
