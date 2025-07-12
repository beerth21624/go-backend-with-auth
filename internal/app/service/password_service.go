package service

import (
	"beerdosan-backend/internal/app/domain"
)

type PasswordService interface {
	Hash(password string) (domain.HashedPassword, error)
	Verify(password string, hashedPassword domain.HashedPassword) error
	ValidateStrength(password string) error
	GenerateRandomPassword(length int) (string, error)
	GenerateSecureToken(length int) (string, error)
}
