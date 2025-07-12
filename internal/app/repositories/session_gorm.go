package repositories

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"venturex-backend/internal/app/domain"
)

type SessionModel struct {
	ID                string `gorm:"type:uuid;primaryKey"`
	UserID            string `gorm:"type:uuid;not null;index"`
	AccessToken       string `gorm:"type:text;not null"`
	RefreshTokenValue string `gorm:"type:varchar(255);not null;uniqueIndex"`
	DeviceFingerprint string `gorm:"type:varchar(500);not null"`
	IPAddress         string `gorm:"type:varchar(45);not null"`
	UserAgent         string `gorm:"type:text;not null"`
	IsActive          bool   `gorm:"default:true"`
	ExpiresAt         time.Time
	RefreshExpiresAt  time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time

	User UserModel `gorm:"foreignKey:UserID"`
}

func (SessionModel) TableName() string {
	return "sessions"
}

func (s *SessionModel) ToDomain() (*domain.Session, error) {
	return domain.ReconstructSession(
		s.ID,
		s.UserID,
		s.AccessToken,
		s.RefreshTokenValue,
		s.DeviceFingerprint,
		s.IPAddress,
		s.UserAgent,
		s.IsActive,
		s.ExpiresAt,
		s.RefreshExpiresAt,
		s.CreatedAt,
		s.UpdatedAt,
		s.UpdatedAt,
	)
}

func CreateSessionModelFromDomain(session *domain.Session) *SessionModel {
	return &SessionModel{
		ID:                session.ID().String(),
		UserID:            session.UserID().String(),
		AccessToken:       session.AccessToken().String(),
		RefreshTokenValue: session.RefreshTokenValue().String(),
		DeviceFingerprint: session.DeviceFingerprint().String(),
		IPAddress:         session.IPAddress().String(),
		UserAgent:         session.UserAgent().String(),
		IsActive:          session.IsActive(),
		ExpiresAt:         session.ExpiresAt().Time(),
		RefreshExpiresAt:  session.RefreshExpiresAt().Time(),
		CreatedAt:         session.CreatedAt().Time(),
		UpdatedAt:         session.UpdatedAt().Time(),
	}
}

func CreateNewSessionModelFromDomain(session *domain.Session) *SessionModel {
	return &SessionModel{
		ID:                session.ID().String(),
		UserID:            session.UserID().String(),
		AccessToken:       session.AccessToken().String(),
		RefreshTokenValue: session.RefreshTokenValue().String(),
		DeviceFingerprint: session.DeviceFingerprint().String(),
		IPAddress:         session.IPAddress().String(),
		UserAgent:         session.UserAgent().String(),
		IsActive:          session.IsActive(),
		ExpiresAt:         session.ExpiresAt().Time(),
		RefreshExpiresAt:  session.RefreshExpiresAt().Time(),
		CreatedAt:         session.CreatedAt().Time(),
		UpdatedAt:         session.UpdatedAt().Time(),
	}
}

func (r *SessionRepositoryGorm) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	return r.CreateInTx(r.db.WithContext(ctx), session)
}

func (r *SessionRepositoryGorm) CreateInTx(tx *gorm.DB, session *domain.Session) (*domain.Session, error) {
	model := CreateNewSessionModelFromDomain(session)

	if err := tx.Create(model).Error; err != nil {
		return nil, err
	}

	return model.ToDomain()
}

func (r *SessionRepositoryGorm) Update(ctx context.Context, session *domain.Session) error {
	return r.UpdateInTx(r.db.WithContext(ctx), session)
}

func (r *SessionRepositoryGorm) UpdateInTx(tx *gorm.DB, session *domain.Session) error {
	model := CreateSessionModelFromDomain(session)
	return tx.Save(model).Error
}

func (r *SessionRepositoryGorm) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	var model SessionModel
	err := r.db.WithContext(ctx).Where("refresh_token_value = ? AND is_active = true", refreshToken).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain()
}

func (r *SessionRepositoryGorm) GetBySessionID(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	var model SessionModel
	err := r.db.WithContext(ctx).Where("id = ?", sessionID.String()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain()
}

func (r *SessionRepositoryGorm) GetActiveSessionsByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Session, error) {
	var models []SessionModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID.String()).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	sessions := make([]*domain.Session, len(models))
	for i, model := range models {
		session, err := model.ToDomain()
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *SessionRepositoryGorm) InvalidateSession(ctx context.Context, sessionID domain.SessionID) error {
	return r.db.WithContext(ctx).Model(&SessionModel{}).
		Where("id = ?", sessionID.String()).
		Update("is_active", false).Error
}

func (r *SessionRepositoryGorm) InvalidateAllUserSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error {
	query := r.db.WithContext(ctx).Model(&SessionModel{}).
		Where("user_id = ? AND is_active = true", userID.String())

	if excludeSessionID.String() != "" {
		query = query.Where("id != ?", excludeSessionID.String())
	}

	return query.Update("is_active", false).Error
}

func (r *SessionRepositoryGorm) UpdateLastActivity(ctx context.Context, sessionID domain.SessionID) error {
	return r.db.WithContext(ctx).Model(&SessionModel{}).
		Where("id = ? AND is_active = true", sessionID.String()).
		Update("updated_at", time.Now()).Error
}
