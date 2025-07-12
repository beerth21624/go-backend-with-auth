package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"venturex-backend/internal/app/domain"
)

type LoginAttemptModel struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	Username      string    `gorm:"type:varchar(255);not null;index"`
	IPAddress     string    `gorm:"type:varchar(45);not null;index"`
	UserAgent     string    `gorm:"type:text;not null"`
	Success       bool      `gorm:"default:false;index"`
	FailureReason *string   `gorm:"type:varchar(255)"`
	AttemptedAt   time.Time `gorm:"index"`
}

func (LoginAttemptModel) TableName() string {
	return "login_attempts"
}

func (la *LoginAttemptModel) ToDomain() (*domain.LoginAttempt, error) {
	return domain.ReconstructLoginAttempt(
		int64(la.ID),
		la.Username,
		la.IPAddress,
		la.UserAgent,
		la.Success,
		la.FailureReason,
		la.AttemptedAt,
	)
}

func CreateLoginAttemptModelFromDomain(attempt *domain.LoginAttempt) *LoginAttemptModel {
	return &LoginAttemptModel{
		ID:            uint64(attempt.ID().Int64()),
		Username:      attempt.Username().String(),
		IPAddress:     attempt.IPAddress().String(),
		UserAgent:     attempt.UserAgent().String(),
		Success:       attempt.Success(),
		FailureReason: attempt.FailureReason(),
		AttemptedAt:   attempt.AttemptedAt().Time(),
	}
}

func CreateNewLoginAttemptModelFromDomain(attempt *domain.LoginAttempt) *LoginAttemptModel {
	return &LoginAttemptModel{
		Username:      attempt.Username().String(),
		IPAddress:     attempt.IPAddress().String(),
		UserAgent:     attempt.UserAgent().String(),
		Success:       attempt.Success(),
		FailureReason: attempt.FailureReason(),
		AttemptedAt:   attempt.AttemptedAt().Time(),
	}
}

func (r *LoginAttemptRepositoryGorm) Create(ctx context.Context, attempt *domain.LoginAttempt) error {
	return r.CreateInTx(r.db.WithContext(ctx), attempt)
}

func (r *LoginAttemptRepositoryGorm) CreateInTx(tx *gorm.DB, attempt *domain.LoginAttempt) error {
	model := CreateNewLoginAttemptModelFromDomain(attempt)

	if err := tx.Create(model).Error; err != nil {
		return err
	}

	return nil
}

func (r *LoginAttemptRepositoryGorm) CountFailedAttemptsByUsernameAndIP(ctx context.Context, username, ipAddress string, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&LoginAttemptModel{}).
		Where("username = ? AND ip_address = ? AND success = false AND attempted_at >= ?", username, ipAddress, since).
		Count(&count).Error
	return count, err
}
