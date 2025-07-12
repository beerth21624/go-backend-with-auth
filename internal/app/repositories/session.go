package repositories

import (
	"context"
	"beerdosan-backend/internal/app/domain"
	"beerdosan-backend/internal/pkg/database"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) (*domain.Session, error)
	CreateInTx(tx *gorm.DB, session *domain.Session) (*domain.Session, error)
	GetBySessionID(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	UpdateInTx(tx *gorm.DB, session *domain.Session) error

	FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	GetActiveSessionsByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Session, error)

	InvalidateSession(ctx context.Context, sessionID domain.SessionID) error
	InvalidateAllUserSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error
	UpdateLastActivity(ctx context.Context, sessionID domain.SessionID) error
}

type SessionRepositoryGorm struct {
	db *database.Database
}

func NewSessionRepository(db *database.Database) *SessionRepositoryGorm {
	return &SessionRepositoryGorm{db: db}
}

var _ SessionRepository = (*SessionRepositoryGorm)(nil)
