package database

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManagerInterface interface {
	ExecuteInTransaction(ctx context.Context, fn func(context.Context) error) error
	Execute(ctx context.Context, fn func(*UnitOfWork) error) error
	ExecuteWithOptions(ctx context.Context, opts *TransactionOptions, fn func(*gorm.DB) error) error
	RetryableTransaction(ctx context.Context, maxRetries int, fn func(*gorm.DB) error) error
	WithTransaction(ctx context.Context) (*TransactionContext, func() error, error)
	NewUnitOfWork(ctx context.Context) *UnitOfWork
	NewBatch(batchSize int) *Batch
	NewSaga() *Saga
}

var _ TransactionManagerInterface = (*TransactionManager)(nil)
