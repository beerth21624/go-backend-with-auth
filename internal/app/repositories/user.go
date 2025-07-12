package repositories

import (
	"context"
	"beerdosan-backend/internal/app/domain"
	"beerdosan-backend/internal/pkg/database"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type UserRepositoryGorm struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepositoryGorm {
	return &UserRepositoryGorm{db: db}
}

var _ UserRepository = (*UserRepositoryGorm)(nil)
