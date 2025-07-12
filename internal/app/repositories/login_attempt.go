package repositories

import (
	"context"
	"time"
	"beerdosan-backend/internal/app/domain"
	"beerdosan-backend/internal/pkg/database"

	"gorm.io/gorm"
)

type LoginAttemptRepository interface {
	Create(ctx context.Context, attempt *domain.LoginAttempt) error
	CreateInTx(tx *gorm.DB, attempt *domain.LoginAttempt) error
	CountFailedAttemptsByUsernameAndIP(ctx context.Context, username, ipAddress string, since time.Time) (int64, error)
}

type LoginAttemptRepositoryGorm struct {
	db *database.Database
}

func NewLoginAttemptRepository(db *database.Database) *LoginAttemptRepositoryGorm {
	return &LoginAttemptRepositoryGorm{db: db}
}

var _ LoginAttemptRepository = (*LoginAttemptRepositoryGorm)(nil)
