package service

import (
	"context"
	"beerdosan-backend/internal/app/domain"
	"beerdosan-backend/internal/app/repositories"
)

type AuthService interface {
	ValidateCredentials(ctx context.Context, username, password string) (*domain.User, error)
	CreateSession(ctx context.Context, userID domain.UserID, deviceInfo, ipAddress string) (*domain.Session, error)
	ValidateSession(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error)
	InvalidateSession(ctx context.Context, sessionID domain.SessionID) error
	InvalidateAllUserSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error
	UpdateSessionActivity(ctx context.Context, sessionID domain.SessionID) error
	RecordLoginAttempt(ctx context.Context, username, ipAddress string, success bool, failureReason string) error
	CheckRateLimit(ctx context.Context, username, ipAddress string) error
	ValidateToken(ctx context.Context, token string) (*AuthClaims, error)
}

type AuthServiceImpl struct {
	userRepo         repositories.UserRepository
	sessionRepo      repositories.SessionRepository
	loginAttemptRepo repositories.LoginAttemptRepository
	passwordService  PasswordService
	jwtService       JWTService
}

func NewAuthService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	loginAttemptRepo repositories.LoginAttemptRepository,
	passwordService PasswordService,
	jwtService JWTService,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:         userRepo,
		sessionRepo:      sessionRepo,
		loginAttemptRepo: loginAttemptRepo,
		passwordService:  passwordService,
		jwtService:       jwtService,
	}
}

var _ AuthService = (*AuthServiceImpl)(nil)
