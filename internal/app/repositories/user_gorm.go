package repositories

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"beerdosan-backend/internal/app/domain"
)

type UserModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Username  string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Email     string `gorm:"type:varchar(255);uniqueIndex;not null"`
	FirstName string `gorm:"type:varchar(100);not null"`
	LastName  string `gorm:"type:varchar(100);not null"`
	Password  string `gorm:"type:text;not null"`
	Role      string `gorm:"type:varchar(20);default:'user'"`
	Status    string `gorm:"type:varchar(20);default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}

func (u *UserModel) ToDomain() (*domain.User, error) {
	return domain.ReconstructUser(
		u.ID,
		u.Username,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
		u.Role,
		u.Status,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

func CreateModelFromDomain(user *domain.User) *UserModel {
	return &UserModel{
		ID:        user.ID().String(),
		Username:  user.Username().String(),
		Email:     user.Email().String(),
		FirstName: user.FirstName().String(),
		LastName:  user.LastName().String(),
		Password:  user.Password().String(),
		Role:      user.Role().String(),
		Status:    user.Status().String(),
		CreatedAt: user.CreatedAt().Time(),
		UpdatedAt: user.UpdatedAt().Time(),
	}
}

func CreateNewModelFromDomain(user *domain.User) *UserModel {
	return &UserModel{
		Username:  user.Username().String(),
		Email:     user.Email().String(),
		FirstName: user.FirstName().String(),
		LastName:  user.LastName().String(),
		Password:  user.Password().String(),
		Role:      user.Role().String(),
		Status:    user.Status().String(),
		CreatedAt: user.CreatedAt().Time(),
		UpdatedAt: user.UpdatedAt().Time(),
	}
}

func (r *UserRepositoryGorm) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := CreateNewModelFromDomain(user)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return model.ToDomain()
}

func (r *UserRepositoryGorm) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain()
}

func (r *UserRepositoryGorm) Update(ctx context.Context, user *domain.User) error {
	model := CreateModelFromDomain(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *UserRepositoryGorm) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain()
}
