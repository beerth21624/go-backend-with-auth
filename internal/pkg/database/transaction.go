package database

import (
	"context"
	"database/sql"
	"fmt"

	"gorm.io/gorm"
)

type TransactionManager struct {
	db *Database
}

func NewTransactionManager(db *Database) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(context.Context) error) error {
	return tm.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

type UnitOfWork struct {
	tx        *gorm.DB
	commits   []func() error
	rollbacks []func() error
}

func (tm *TransactionManager) NewUnitOfWork(ctx context.Context) *UnitOfWork {
	return &UnitOfWork{
		commits:   make([]func() error, 0),
		rollbacks: make([]func() error, 0),
	}
}

func (tm *TransactionManager) Execute(ctx context.Context, fn func(*UnitOfWork) error) error {
	uow := tm.NewUnitOfWork(ctx)

	return tm.db.Transaction(ctx, func(tx *gorm.DB) error {
		uow.tx = tx

		if err := fn(uow); err != nil {
			for _, rollback := range uow.rollbacks {
				if rollbackErr := rollback(); rollbackErr != nil {
					fmt.Printf("Rollback error: %v\n", rollbackErr)
				}
			}
			return err
		}

		for _, commit := range uow.commits {
			if err := commit(); err != nil {
				return fmt.Errorf("commit hook failed: %w", err)
			}
		}

		return nil
	})
}

func (uow *UnitOfWork) GetTx() *gorm.DB {
	return uow.tx
}

func (uow *UnitOfWork) OnCommit(fn func() error) {
	uow.commits = append(uow.commits, fn)
}

func (uow *UnitOfWork) OnRollback(fn func() error) {
	uow.rollbacks = append(uow.rollbacks, fn)
}

type TransactionOptions struct {
	IsolationLevel sql.IsolationLevel
	ReadOnly       bool
}

func WithReadCommitted() *TransactionOptions {
	return &TransactionOptions{
		IsolationLevel: sql.LevelReadCommitted,
	}
}

func WithReadUncommitted() *TransactionOptions {
	return &TransactionOptions{
		IsolationLevel: sql.LevelReadUncommitted,
	}
}

func WithRepeatableRead() *TransactionOptions {
	return &TransactionOptions{
		IsolationLevel: sql.LevelRepeatableRead,
	}
}

func WithSerializable() *TransactionOptions {
	return &TransactionOptions{
		IsolationLevel: sql.LevelSerializable,
	}
}

func (opt *TransactionOptions) WithReadOnly() *TransactionOptions {
	opt.ReadOnly = true
	return opt
}

func (tm *TransactionManager) ExecuteWithOptions(ctx context.Context, opts *TransactionOptions, fn func(*gorm.DB) error) error {
	if opts == nil {
		return tm.db.Transaction(ctx, fn)
	}

	sqlOpts := &sql.TxOptions{
		Isolation: opts.IsolationLevel,
		ReadOnly:  opts.ReadOnly,
	}

	return tm.db.TransactionWithOptions(ctx, sqlOpts, fn)
}

type Batch struct {
	operations []func(*gorm.DB) error
	batchSize  int
}

func (tm *TransactionManager) NewBatch(batchSize int) *Batch {
	if batchSize <= 0 {
		batchSize = 100
	}
	return &Batch{
		operations: make([]func(*gorm.DB) error, 0),
		batchSize:  batchSize,
	}
}

func (b *Batch) Add(op func(*gorm.DB) error) {
	b.operations = append(b.operations, op)
}

func (b *Batch) Execute(ctx context.Context, tm *TransactionManager) error {
	if len(b.operations) == 0 {
		return nil
	}

	for i := 0; i < len(b.operations); i += b.batchSize {
		end := i + b.batchSize
		if end > len(b.operations) {
			end = len(b.operations)
		}

		batch := b.operations[i:end]

		err := tm.Execute(ctx, func(uow *UnitOfWork) error {
			for _, op := range batch {
				if err := op(uow.GetTx()); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("batch execution failed at index %d: %w", i, err)
		}
	}

	return nil
}

type Saga struct {
	steps []SagaStep
}

type SagaStep struct {
	Execute    func(*gorm.DB) error
	Compensate func(*gorm.DB) error
}

func (tm *TransactionManager) NewSaga() *Saga {
	return &Saga{
		steps: make([]SagaStep, 0),
	}
}

func (s *Saga) AddStep(execute, compensate func(*gorm.DB) error) {
	s.steps = append(s.steps, SagaStep{
		Execute:    execute,
		Compensate: compensate,
	})
}

func (s *Saga) Execute(ctx context.Context, tm *TransactionManager) error {
	executedSteps := make([]int, 0)

	for i, step := range s.steps {
		err := tm.Execute(ctx, func(uow *UnitOfWork) error {
			return step.Execute(uow.GetTx())
		})

		if err != nil {
			for j := len(executedSteps) - 1; j >= 0; j-- {
				stepIndex := executedSteps[j]
				compensateErr := tm.Execute(ctx, func(uow *UnitOfWork) error {
					return s.steps[stepIndex].Compensate(uow.GetTx())
				})

				if compensateErr != nil {
					return fmt.Errorf("saga compensation failed at step %d: %w (original error: %v)", stepIndex, compensateErr, err)
				}
			}

			return fmt.Errorf("saga failed at step %d: %w", i, err)
		}

		executedSteps = append(executedSteps, i)
	}

	return nil
}

func (tm *TransactionManager) RetryableTransaction(ctx context.Context, maxRetries int, fn func(*gorm.DB) error) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := tm.db.Transaction(ctx, fn)
		if err == nil {
			return nil
		}

		if !isRetryableError(err) {
			return err
		}

		lastErr = err

		if attempt == maxRetries {
			break
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return fmt.Errorf("transaction failed after %d retries: %w", maxRetries, lastErr)
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	retryablePatterns := []string{
		"deadlock",
		"serialization failure",
		"could not serialize access",
		"retry transaction",
	}

	for _, pattern := range retryablePatterns {
		if containsSubstring(errStr, pattern) {
			return true
		}
	}

	return false
}

type TransactionContext struct {
	ctx context.Context
	tx  *gorm.DB
}

func (tm *TransactionManager) WithTransaction(ctx context.Context) (*TransactionContext, func() error, error) {
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	txCtx := &TransactionContext{
		ctx: ctx,
		tx:  tx,
	}

	cleanup := func() error {
		return tx.Rollback().Error
	}

	return txCtx, cleanup, nil
}

func (tc *TransactionContext) Commit() error {
	return tc.tx.Commit().Error
}

func (tc *TransactionContext) Rollback() error {
	return tc.tx.Rollback().Error
}

func (tc *TransactionContext) DB() *gorm.DB {
	return tc.tx
}

func (tc *TransactionContext) Context() context.Context {
	return tc.ctx
}
