package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	CreateInTx(tx *gorm.DB, entity *T) error

	GetByID(ctx context.Context, id interface{}) (*T, error)
	GetByIDInTx(tx *gorm.DB, id interface{}) (*T, error)

	Update(ctx context.Context, entity *T) error
	UpdateInTx(tx *gorm.DB, entity *T) error

	Delete(ctx context.Context, id interface{}) error
	DeleteInTx(tx *gorm.DB, id interface{}) error

	List(ctx context.Context, query *Query) ([]*T, error)
	ListInTx(tx *gorm.DB, query *Query) ([]*T, error)

	Count(ctx context.Context, query *Query) (int64, error)
	CountInTx(tx *gorm.DB, query *Query) (int64, error)

	Exists(ctx context.Context, id interface{}) (bool, error)
	ExistsInTx(tx *gorm.DB, id interface{}) (bool, error)
}

type BaseRepository[T any] struct {
	db    *Database
	model T
}

func NewBaseRepository[T any](db *Database) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:    db,
		model: *new(T),
	}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) CreateInTx(tx *gorm.DB, entity *T) error {
	return tx.Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) GetByIDInTx(tx *gorm.DB, id interface{}) (*T, error) {
	var entity T
	err := tx.First(&entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) UpdateInTx(tx *gorm.DB, entity *T) error {
	return tx.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id interface{}) error {
	return r.db.WithContext(ctx).Delete(&r.model, id).Error
}

func (r *BaseRepository[T]) DeleteInTx(tx *gorm.DB, id interface{}) error {
	return tx.Delete(&r.model, id).Error
}

func (r *BaseRepository[T]) List(ctx context.Context, query *Query) ([]*T, error) {
	return r.ListInTx(r.db.WithContext(ctx), query)
}

func (r *BaseRepository[T]) ListInTx(tx *gorm.DB, query *Query) ([]*T, error) {
	var entities []*T
	db := tx.Model(&r.model)

	if query != nil {
		db = r.applyQuery(db, query)
	}

	err := db.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Count(ctx context.Context, query *Query) (int64, error) {
	return r.CountInTx(r.db.WithContext(ctx), query)
}

func (r *BaseRepository[T]) CountInTx(tx *gorm.DB, query *Query) (int64, error) {
	var count int64
	db := tx.Model(&r.model)

	if query != nil {
		db = r.applyQuery(db, query)
	}

	err := db.Count(&count).Error
	return count, err
}

func (r *BaseRepository[T]) Exists(ctx context.Context, id interface{}) (bool, error) {
	return r.ExistsInTx(r.db.WithContext(ctx), id)
}

func (r *BaseRepository[T]) ExistsInTx(tx *gorm.DB, id interface{}) (bool, error) {
	var count int64
	err := tx.Model(&r.model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *BaseRepository[T]) applyQuery(db *gorm.DB, query *Query) *gorm.DB {
	if query == nil {
		return db
	}

	for _, filter := range query.Filters {
		db = db.Where(filter.Field+" "+filter.Operator+" ?", filter.Value)
	}

	for _, sort := range query.Sorts {
		direction := "ASC"
		if sort.Desc {
			direction = "DESC"
		}
		db = db.Order(fmt.Sprintf("%s %s", sort.Field, direction))
	}

	if query.Limit > 0 {
		db = db.Limit(int(query.Limit))
	}
	if query.Offset > 0 {
		db = db.Offset(int(query.Offset))
	}

	for _, preload := range query.Preloads {
		if preload.Conditions != nil {
			db = db.Preload(preload.Field, preload.Conditions...)
		} else {
			db = db.Preload(preload.Field)
		}
	}

	return db
}

func (r *BaseRepository[T]) GetDB() *Database {
	return r.db
}

func (r *BaseRepository[T]) GetModel() T {
	return r.model
}

func (r *BaseRepository[T]) GetModelName() string {
	return reflect.TypeOf(r.model).Name()
}
