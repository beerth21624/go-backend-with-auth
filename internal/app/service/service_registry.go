package service

import (
	"beerdosan-backend/internal/app/repositories"
	"beerdosan-backend/internal/pkg/jwt"
	"beerdosan-backend/internal/pkg/password"
)

type ServiceRegistry struct {
	authService     AuthService
	jwtService      JWTService
	passwordService PasswordService
}

func NewServiceRegistry(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	loginAttemptRepo repositories.LoginAttemptRepository,
	jwtService jwt.JWTService,
	passwordService password.PasswordService,
) *ServiceRegistry {
	pwdService := NewPasswordService(passwordService)

	jwtSvc := NewJWTService(jwtService)

	authSvc := NewAuthService(
		userRepo,
		sessionRepo,
		loginAttemptRepo,
		pwdService,
		jwtSvc,
	)

	return &ServiceRegistry{
		authService:     authSvc,
		jwtService:      jwtSvc,
		passwordService: pwdService,
	}
}

func (r *ServiceRegistry) AuthService() AuthService {
	return r.authService
}

func (r *ServiceRegistry) JWTService() JWTService {
	return r.jwtService
}

func (r *ServiceRegistry) PasswordService() PasswordService {
	return r.passwordService
}
