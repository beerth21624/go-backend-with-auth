// Code generated by mockery v2.53.4. DO NOT EDIT.

package database

import (
	context "context"
	database "beerdosan-backend/internal/pkg/database"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// MockTransactionManagerInterface is an autogenerated mock type for the TransactionManagerInterface type
type MockTransactionManagerInterface struct {
	mock.Mock
}

type MockTransactionManagerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTransactionManagerInterface) EXPECT() *MockTransactionManagerInterface_Expecter {
	return &MockTransactionManagerInterface_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, fn
func (_m *MockTransactionManagerInterface) Execute(ctx context.Context, fn func(*database.UnitOfWork) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(*database.UnitOfWork) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransactionManagerInterface_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockTransactionManagerInterface_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(*database.UnitOfWork) error
func (_e *MockTransactionManagerInterface_Expecter) Execute(ctx interface{}, fn interface{}) *MockTransactionManagerInterface_Execute_Call {
	return &MockTransactionManagerInterface_Execute_Call{Call: _e.mock.On("Execute", ctx, fn)}
}

func (_c *MockTransactionManagerInterface_Execute_Call) Run(run func(ctx context.Context, fn func(*database.UnitOfWork) error)) *MockTransactionManagerInterface_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(*database.UnitOfWork) error))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_Execute_Call) Return(_a0 error) *MockTransactionManagerInterface_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_Execute_Call) RunAndReturn(run func(context.Context, func(*database.UnitOfWork) error) error) *MockTransactionManagerInterface_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteInTransaction provides a mock function with given fields: ctx, fn
func (_m *MockTransactionManagerInterface) ExecuteInTransaction(ctx context.Context, fn func(context.Context) error) error {
	ret := _m.Called(ctx, fn)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteInTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransactionManagerInterface_ExecuteInTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteInTransaction'
type MockTransactionManagerInterface_ExecuteInTransaction_Call struct {
	*mock.Call
}

// ExecuteInTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(context.Context) error
func (_e *MockTransactionManagerInterface_Expecter) ExecuteInTransaction(ctx interface{}, fn interface{}) *MockTransactionManagerInterface_ExecuteInTransaction_Call {
	return &MockTransactionManagerInterface_ExecuteInTransaction_Call{Call: _e.mock.On("ExecuteInTransaction", ctx, fn)}
}

func (_c *MockTransactionManagerInterface_ExecuteInTransaction_Call) Run(run func(ctx context.Context, fn func(context.Context) error)) *MockTransactionManagerInterface_ExecuteInTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_ExecuteInTransaction_Call) Return(_a0 error) *MockTransactionManagerInterface_ExecuteInTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_ExecuteInTransaction_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *MockTransactionManagerInterface_ExecuteInTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteWithOptions provides a mock function with given fields: ctx, opts, fn
func (_m *MockTransactionManagerInterface) ExecuteWithOptions(ctx context.Context, opts *database.TransactionOptions, fn func(*gorm.DB) error) error {
	ret := _m.Called(ctx, opts, fn)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteWithOptions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *database.TransactionOptions, func(*gorm.DB) error) error); ok {
		r0 = rf(ctx, opts, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransactionManagerInterface_ExecuteWithOptions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteWithOptions'
type MockTransactionManagerInterface_ExecuteWithOptions_Call struct {
	*mock.Call
}

// ExecuteWithOptions is a helper method to define mock.On call
//   - ctx context.Context
//   - opts *database.TransactionOptions
//   - fn func(*gorm.DB) error
func (_e *MockTransactionManagerInterface_Expecter) ExecuteWithOptions(ctx interface{}, opts interface{}, fn interface{}) *MockTransactionManagerInterface_ExecuteWithOptions_Call {
	return &MockTransactionManagerInterface_ExecuteWithOptions_Call{Call: _e.mock.On("ExecuteWithOptions", ctx, opts, fn)}
}

func (_c *MockTransactionManagerInterface_ExecuteWithOptions_Call) Run(run func(ctx context.Context, opts *database.TransactionOptions, fn func(*gorm.DB) error)) *MockTransactionManagerInterface_ExecuteWithOptions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*database.TransactionOptions), args[2].(func(*gorm.DB) error))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_ExecuteWithOptions_Call) Return(_a0 error) *MockTransactionManagerInterface_ExecuteWithOptions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_ExecuteWithOptions_Call) RunAndReturn(run func(context.Context, *database.TransactionOptions, func(*gorm.DB) error) error) *MockTransactionManagerInterface_ExecuteWithOptions_Call {
	_c.Call.Return(run)
	return _c
}

// NewBatch provides a mock function with given fields: batchSize
func (_m *MockTransactionManagerInterface) NewBatch(batchSize int) *database.Batch {
	ret := _m.Called(batchSize)

	if len(ret) == 0 {
		panic("no return value specified for NewBatch")
	}

	var r0 *database.Batch
	if rf, ok := ret.Get(0).(func(int) *database.Batch); ok {
		r0 = rf(batchSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.Batch)
		}
	}

	return r0
}

// MockTransactionManagerInterface_NewBatch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewBatch'
type MockTransactionManagerInterface_NewBatch_Call struct {
	*mock.Call
}

// NewBatch is a helper method to define mock.On call
//   - batchSize int
func (_e *MockTransactionManagerInterface_Expecter) NewBatch(batchSize interface{}) *MockTransactionManagerInterface_NewBatch_Call {
	return &MockTransactionManagerInterface_NewBatch_Call{Call: _e.mock.On("NewBatch", batchSize)}
}

func (_c *MockTransactionManagerInterface_NewBatch_Call) Run(run func(batchSize int)) *MockTransactionManagerInterface_NewBatch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_NewBatch_Call) Return(_a0 *database.Batch) *MockTransactionManagerInterface_NewBatch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_NewBatch_Call) RunAndReturn(run func(int) *database.Batch) *MockTransactionManagerInterface_NewBatch_Call {
	_c.Call.Return(run)
	return _c
}

// NewSaga provides a mock function with no fields
func (_m *MockTransactionManagerInterface) NewSaga() *database.Saga {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NewSaga")
	}

	var r0 *database.Saga
	if rf, ok := ret.Get(0).(func() *database.Saga); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.Saga)
		}
	}

	return r0
}

// MockTransactionManagerInterface_NewSaga_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewSaga'
type MockTransactionManagerInterface_NewSaga_Call struct {
	*mock.Call
}

// NewSaga is a helper method to define mock.On call
func (_e *MockTransactionManagerInterface_Expecter) NewSaga() *MockTransactionManagerInterface_NewSaga_Call {
	return &MockTransactionManagerInterface_NewSaga_Call{Call: _e.mock.On("NewSaga")}
}

func (_c *MockTransactionManagerInterface_NewSaga_Call) Run(run func()) *MockTransactionManagerInterface_NewSaga_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTransactionManagerInterface_NewSaga_Call) Return(_a0 *database.Saga) *MockTransactionManagerInterface_NewSaga_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_NewSaga_Call) RunAndReturn(run func() *database.Saga) *MockTransactionManagerInterface_NewSaga_Call {
	_c.Call.Return(run)
	return _c
}

// NewUnitOfWork provides a mock function with given fields: ctx
func (_m *MockTransactionManagerInterface) NewUnitOfWork(ctx context.Context) *database.UnitOfWork {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for NewUnitOfWork")
	}

	var r0 *database.UnitOfWork
	if rf, ok := ret.Get(0).(func(context.Context) *database.UnitOfWork); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.UnitOfWork)
		}
	}

	return r0
}

// MockTransactionManagerInterface_NewUnitOfWork_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewUnitOfWork'
type MockTransactionManagerInterface_NewUnitOfWork_Call struct {
	*mock.Call
}

// NewUnitOfWork is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockTransactionManagerInterface_Expecter) NewUnitOfWork(ctx interface{}) *MockTransactionManagerInterface_NewUnitOfWork_Call {
	return &MockTransactionManagerInterface_NewUnitOfWork_Call{Call: _e.mock.On("NewUnitOfWork", ctx)}
}

func (_c *MockTransactionManagerInterface_NewUnitOfWork_Call) Run(run func(ctx context.Context)) *MockTransactionManagerInterface_NewUnitOfWork_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_NewUnitOfWork_Call) Return(_a0 *database.UnitOfWork) *MockTransactionManagerInterface_NewUnitOfWork_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_NewUnitOfWork_Call) RunAndReturn(run func(context.Context) *database.UnitOfWork) *MockTransactionManagerInterface_NewUnitOfWork_Call {
	_c.Call.Return(run)
	return _c
}

// RetryableTransaction provides a mock function with given fields: ctx, maxRetries, fn
func (_m *MockTransactionManagerInterface) RetryableTransaction(ctx context.Context, maxRetries int, fn func(*gorm.DB) error) error {
	ret := _m.Called(ctx, maxRetries, fn)

	if len(ret) == 0 {
		panic("no return value specified for RetryableTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, func(*gorm.DB) error) error); ok {
		r0 = rf(ctx, maxRetries, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockTransactionManagerInterface_RetryableTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RetryableTransaction'
type MockTransactionManagerInterface_RetryableTransaction_Call struct {
	*mock.Call
}

// RetryableTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - maxRetries int
//   - fn func(*gorm.DB) error
func (_e *MockTransactionManagerInterface_Expecter) RetryableTransaction(ctx interface{}, maxRetries interface{}, fn interface{}) *MockTransactionManagerInterface_RetryableTransaction_Call {
	return &MockTransactionManagerInterface_RetryableTransaction_Call{Call: _e.mock.On("RetryableTransaction", ctx, maxRetries, fn)}
}

func (_c *MockTransactionManagerInterface_RetryableTransaction_Call) Run(run func(ctx context.Context, maxRetries int, fn func(*gorm.DB) error)) *MockTransactionManagerInterface_RetryableTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(func(*gorm.DB) error))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_RetryableTransaction_Call) Return(_a0 error) *MockTransactionManagerInterface_RetryableTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTransactionManagerInterface_RetryableTransaction_Call) RunAndReturn(run func(context.Context, int, func(*gorm.DB) error) error) *MockTransactionManagerInterface_RetryableTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// WithTransaction provides a mock function with given fields: ctx
func (_m *MockTransactionManagerInterface) WithTransaction(ctx context.Context) (*database.TransactionContext, func() error, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for WithTransaction")
	}

	var r0 *database.TransactionContext
	var r1 func() error
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) (*database.TransactionContext, func() error, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *database.TransactionContext); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*database.TransactionContext)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) func() error); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func() error)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockTransactionManagerInterface_WithTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTransaction'
type MockTransactionManagerInterface_WithTransaction_Call struct {
	*mock.Call
}

// WithTransaction is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockTransactionManagerInterface_Expecter) WithTransaction(ctx interface{}) *MockTransactionManagerInterface_WithTransaction_Call {
	return &MockTransactionManagerInterface_WithTransaction_Call{Call: _e.mock.On("WithTransaction", ctx)}
}

func (_c *MockTransactionManagerInterface_WithTransaction_Call) Run(run func(ctx context.Context)) *MockTransactionManagerInterface_WithTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockTransactionManagerInterface_WithTransaction_Call) Return(_a0 *database.TransactionContext, _a1 func() error, _a2 error) *MockTransactionManagerInterface_WithTransaction_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockTransactionManagerInterface_WithTransaction_Call) RunAndReturn(run func(context.Context) (*database.TransactionContext, func() error, error)) *MockTransactionManagerInterface_WithTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTransactionManagerInterface creates a new instance of MockTransactionManagerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTransactionManagerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTransactionManagerInterface {
	mock := &MockTransactionManagerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
