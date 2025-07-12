package usecase

import (
	"context"

	"beerdosan-backend/internal/app/domain"
	"beerdosan-backend/internal/app/repositories"
	"beerdosan-backend/internal/app/service"
	"beerdosan-backend/internal/pkg/database"
)

type AuthUseCase interface {
	Login(ctx context.Context, req LoginInput) (*LoginOutput, error)
	Logout(ctx context.Context, userID domain.UserID, sessionID domain.SessionID) error
	RefreshToken(ctx context.Context, req RefreshTokenInput) (*RefreshTokenOutput, error)
	GetUserProfile(ctx context.Context, userID domain.UserID) (*GetUserProfileOutput, error)
	GetUserSessions(ctx context.Context, userID domain.UserID) ([]GetUserSessionsOutput, error)
	RevokeSession(ctx context.Context, userID domain.UserID, sessionID domain.SessionID) error
	RevokeAllSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error
	ChangePassword(ctx context.Context, userID domain.UserID, oldPassword, newPassword string) error
}
type AuthUseCaseImpl struct {
	authService     service.AuthService
	jwtService      service.JWTService
	passwordService service.PasswordService
	userRepo        repositories.UserRepository
	sessionRepo     repositories.SessionRepository
	transactionMgr  *database.TransactionManager
}

func NewAuthUseCase(
	authService service.AuthService,
	jwtService service.JWTService,
	passwordService service.PasswordService,
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	transactionMgr *database.TransactionManager,
) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		authService:     authService,
		jwtService:      jwtService,
		passwordService: passwordService,
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		transactionMgr:  transactionMgr,
	}
}

var _ AuthUseCase = (*AuthUseCaseImpl)(nil)
