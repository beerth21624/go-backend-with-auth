package service

import (
	"encoding/base64"
	"fmt"

	"venturex-backend/internal/app/domain"
	"venturex-backend/internal/pkg/password"
)

type passwordServiceImpl struct {
	passwordService password.PasswordService
}

func NewPasswordService(passwordService password.PasswordService) PasswordService {
	return &passwordServiceImpl{
		passwordService: passwordService,
	}
}

func (s *passwordServiceImpl) Hash(password string) (domain.HashedPassword, error) {
	hashed, err := s.passwordService.HashPassword(password)
	if err != nil {
		return "", err
	}
	return domain.HashedPassword(hashed), nil
}

func (s *passwordServiceImpl) Verify(plainPassword string, hashedPassword domain.HashedPassword) error {
	isValid := s.passwordService.VerifyPassword(password.HashedPassword(hashedPassword), plainPassword)
	if !isValid {
		return fmt.Errorf("password verification failed")
	}
	return nil
}

func (s *passwordServiceImpl) ValidateStrength(password string) error {
	return s.passwordService.ValidatePassword(password)
}

func (s *passwordServiceImpl) GenerateRandomPassword(length int) (string, error) {
	return s.passwordService.GenerateRandomPassword(length, true)
}

func (s *passwordServiceImpl) GenerateSecureToken(length int) (string, error) {
	if length < 16 {
		return "", fmt.Errorf("token length must be at least 16")
	}

	password, err := s.passwordService.GenerateRandomPassword(length, false)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString([]byte(password)), nil
}
