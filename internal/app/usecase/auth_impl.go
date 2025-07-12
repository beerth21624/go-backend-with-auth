package usecase

import (
	"context"
	"time"

	"venturex-backend/internal/app/domain"
	"venturex-backend/internal/pkg/sliceutil"
)

type UserInfo struct {
	ID       domain.UserID         `json:"id"`
	Username domain.NonEmptyString `json:"username"`
	Email    domain.NonEmptyString `json:"email"`
	Role     domain.UserRole       `json:"role"`
	Status   domain.Status         `json:"status"`
}

type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
	IPAddress  string `json:"ip_address"`
	RememberMe bool   `json:"remember_me"`
}

type LoginOutput struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

func (uc *AuthUseCaseImpl) Login(ctx context.Context, req LoginInput) (*LoginOutput, error) {
	if err := uc.authService.CheckRateLimit(ctx, req.Username, req.IPAddress); err != nil {
		_ = uc.authService.RecordLoginAttempt(ctx, req.Username, req.IPAddress, false, "rate_limited")
		return nil, err
	}

	user, err := uc.authService.ValidateCredentials(ctx, req.Username, req.Password)
	if err != nil {
		_ = uc.authService.RecordLoginAttempt(ctx, req.Username, req.IPAddress, false, "invalid_credentials")
		return nil, err
	}

	if !user.CanLogin() {
		_ = uc.authService.RecordLoginAttempt(ctx, req.Username, req.IPAddress, false, "account_disabled")
		return nil, domain.ErrAccountLocked
	}

	var response *LoginOutput
	err = uc.transactionMgr.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		session, err := uc.authService.CreateSession(ctx, user.ID(), req.DeviceInfo, req.IPAddress)
		if err != nil {
			return domain.DefineError(domain.ErrCatSystem, "SESSION_CREATE_FAILED", "failed to create session").Wrap(err)
		}

		accessToken, err := uc.jwtService.GenerateAccessToken(user.ID(), session.ID(), user.Role())
		if err != nil {
			return domain.DefineError(domain.ErrCatSystem, "TOKEN_GENERATION_FAILED", "failed to generate access token").Wrap(err)
		}

		refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID(), session.ID())
		if err != nil {
			return domain.DefineError(domain.ErrCatSystem, "TOKEN_GENERATION_FAILED", "failed to generate refresh token").Wrap(err)
		}

		if err := uc.authService.RecordLoginAttempt(ctx, req.Username, req.IPAddress, true, ""); err != nil {
			// Log error but don't fail the login
			// In production, you might want to use a proper logger
			// TODO: Use proper logger
		}

		response = &LoginOutput{
			AccessToken:  accessToken.String(),
			RefreshToken: refreshToken.String(),
			ExpiresAt:    time.Now().Add(15 * time.Minute), // Access token expires in 15 minutes
			User: UserInfo{
				ID:       user.ID(),
				Username: user.Username(),
				Email:    user.Email(),
				Role:     user.Role(),
				Status:   user.Status(),
			},
		}

		return nil
	})

	if err != nil {
		return nil, domain.DefineError(domain.ErrCatSystem, "LOGIN_FAILED", "login process failed").Wrap(err)
	}

	return response, nil
}

func (uc *AuthUseCaseImpl) Logout(ctx context.Context, userID domain.UserID, sessionID domain.SessionID) error {
	session, err := uc.authService.ValidateSession(ctx, sessionID)
	if err != nil {
		return domain.ErrInvalidSession.Wrap(err)
	}

	if session.UserID() != userID {
		return domain.DefineError(domain.ErrCatAuth, "SESSION_USER_MISMATCH", "session does not belong to user")
	}

	if err := uc.authService.InvalidateSession(ctx, sessionID); err != nil {
		return domain.DefineError(domain.ErrCatSystem, "SESSION_INVALIDATE_FAILED", "failed to invalidate session").Wrap(err)
	}

	return nil
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
	DeviceInfo   string `json:"device_info"`
	IPAddress    string `json:"ip_address"`
}

type RefreshTokenOutput struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

func (uc *AuthUseCaseImpl) RefreshToken(ctx context.Context, req RefreshTokenInput) (*RefreshTokenOutput, error) {
	refreshToken, err := domain.NewJWT(req.RefreshToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid.Wrap(err)
	}

	claims, err := uc.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, domain.ErrTokenInvalid.Wrap(err)
	}

	if !claims.IsRefreshToken() {
		return nil, domain.DefineError(domain.ErrCatAuth, "INVALID_TOKEN_TYPE", "token is not a refresh token")
	}

	sessionID, err := domain.NewSessionIDFromString(claims.SessionID)
	if err != nil {
		return nil, domain.ErrInvalidSession.Wrap(err)
	}

	session, err := uc.sessionRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, domain.ErrSessionNotFound.Wrap(err)
	}

	if session == nil || !session.CanRefresh() {
		return nil, domain.ErrInvalidSession
	}

	user, err := uc.userRepo.GetByID(ctx, session.UserID())
	if err != nil {
		return nil, domain.ErrUserNotFound.Wrap(err)
	}

	if !user.CanLogin() {
		return nil, domain.ErrAccountLocked
	}

	newAccessToken, err := uc.jwtService.GenerateAccessToken(user.ID(), session.ID(), user.Role())
	if err != nil {
		return nil, domain.DefineError(domain.ErrCatSystem, "TOKEN_GENERATION_FAILED", "failed to generate access token").Wrap(err)
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID(), session.ID())
	if err != nil {
		return nil, domain.DefineError(domain.ErrCatSystem, "TOKEN_GENERATION_FAILED", "failed to generate refresh token").Wrap(err)
	}

	if err := uc.authService.UpdateSessionActivity(ctx, session.ID()); err != nil {
		// Log error but don't fail the token refresh
		// TODO: Use proper logger
	}

	return &RefreshTokenOutput{
		AccessToken:  newAccessToken.String(),
		RefreshToken: newRefreshToken.String(),
		ExpiresAt:    time.Now().Add(15 * time.Minute),
		User: UserInfo{
			ID:       user.ID(),
			Username: user.Username(),
			Email:    user.Email(),
			Role:     user.Role(),
			Status:   user.Status(),
		},
	}, nil
}

type GetUserProfileOutput struct {
	User UserInfo `json:"user"`
}

func (uc *AuthUseCaseImpl) GetUserProfile(ctx context.Context, userID domain.UserID) (*GetUserProfileOutput, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound.Wrap(err)
	}

	return &GetUserProfileOutput{
		User: UserInfo{
			ID:       user.ID(),
			Username: user.Username(),
			Email:    user.Email(),
			Role:     user.Role(),
			Status:   user.Status(),
		},
	}, nil
}

type GetUserSessionsOutput struct {
	ID           domain.SessionID `json:"id"`
	UserID       domain.UserID    `json:"user_id"`
	DeviceInfo   string           `json:"device_info"`
	IPAddress    string           `json:"ip_address"`
	LastActivity time.Time        `json:"last_activity"`
	CreatedAt    time.Time        `json:"created_at"`
	IsActive     bool             `json:"is_active"`
}

func (uc *AuthUseCaseImpl) GetUserSessions(ctx context.Context, userID domain.UserID) ([]GetUserSessionsOutput, error) {
	sessions, err := uc.sessionRepo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, domain.DefineError(domain.ErrCatSystem, "SESSION_FETCH_FAILED", "failed to get user sessions").Wrap(err)
	}

		sessionInfos := sliceutil.Map(sessions, func(session *domain.Session) GetUserSessionsOutput {
		return GetUserSessionsOutput{
			ID:           session.ID(),
			UserID:       session.UserID(),
			DeviceInfo:   session.DeviceInfo(),
			IPAddress:    session.IPAddress().String(),
			LastActivity: session.LastActivity().Time(),
			CreatedAt:    session.CreatedAt().Time(),
			IsActive:     session.IsActive(),
		}
	})

	return sessionInfos, nil
}

func (uc *AuthUseCaseImpl) RevokeSession(ctx context.Context, userID domain.UserID, sessionID domain.SessionID) error {
	session, err := uc.authService.ValidateSession(ctx, sessionID)
	if err != nil {
		return domain.ErrInvalidSession.Wrap(err)
	}

	if session.UserID() != userID {
		return domain.DefineError(domain.ErrCatAuth, "SESSION_USER_MISMATCH", "session does not belong to user")
	}

	if err := uc.authService.InvalidateSession(ctx, sessionID); err != nil {
		return domain.DefineError(domain.ErrCatSystem, "SESSION_REVOKE_FAILED", "failed to revoke session").Wrap(err)
	}

	return nil
}

func (uc *AuthUseCaseImpl) RevokeAllSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error {
	if err := uc.authService.InvalidateAllUserSessions(ctx, userID, excludeSessionID); err != nil {
		return domain.DefineError(domain.ErrCatSystem, "SESSION_REVOKE_FAILED", "failed to revoke all sessions").Wrap(err)
	}

	return nil
}

func (uc *AuthUseCaseImpl) ChangePassword(ctx context.Context, userID domain.UserID, oldPassword, newPassword string) error {
	if err := uc.passwordService.ValidateStrength(newPassword); err != nil {
		return domain.DefineError(domain.ErrCatValidation, "INVALID_PASSWORD", "password validation failed").Wrap(err)
	}

	return uc.transactionMgr.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		user, err := uc.userRepo.GetByID(ctx, userID)
		if err != nil {
			return domain.ErrUserNotFound.Wrap(err)
		}

		if !user.VerifyPassword(oldPassword) {
			return domain.DefineError(domain.ErrCatAuth, "INVALID_PASSWORD", "old password is incorrect")
		}

		if err := user.ChangePassword(newPassword); err != nil {
			return domain.DefineError(domain.ErrCatSystem, "PASSWORD_CHANGE_FAILED", "failed to change password").Wrap(err)
		}

		if err := uc.userRepo.Update(ctx, user); err != nil {
			return domain.DefineError(domain.ErrCatSystem, "USER_UPDATE_FAILED", "failed to save user").Wrap(err)
		}

		if err := uc.authService.InvalidateAllUserSessions(ctx, userID, domain.SessionID("")); err != nil {
			// Log error but don't fail the password change
			// TODO: Use proper logger
		}

		return nil
	})
}
