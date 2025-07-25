// Code generated by mockery v2.53.4. DO NOT EDIT.

package repositories

import (
	context "context"
	domain "beerdosan-backend/internal/app/domain"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockLoginAttemptRepository is an autogenerated mock type for the LoginAttemptRepository type
type MockLoginAttemptRepository struct {
	mock.Mock
}

type MockLoginAttemptRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLoginAttemptRepository) EXPECT() *MockLoginAttemptRepository_Expecter {
	return &MockLoginAttemptRepository_Expecter{mock: &_m.Mock}
}

// CountFailedAttemptsByUsernameAndIP provides a mock function with given fields: ctx, username, ipAddress, since
func (_m *MockLoginAttemptRepository) CountFailedAttemptsByUsernameAndIP(ctx context.Context, username string, ipAddress string, since time.Time) (int64, error) {
	ret := _m.Called(ctx, username, ipAddress, since)

	if len(ret) == 0 {
		panic("no return value specified for CountFailedAttemptsByUsernameAndIP")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) (int64, error)); ok {
		return rf(ctx, username, ipAddress, since)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) int64); ok {
		r0 = rf(ctx, username, ipAddress, since)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, time.Time) error); ok {
		r1 = rf(ctx, username, ipAddress, since)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountFailedAttemptsByUsernameAndIP'
type MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call struct {
	*mock.Call
}

// CountFailedAttemptsByUsernameAndIP is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - ipAddress string
//   - since time.Time
func (_e *MockLoginAttemptRepository_Expecter) CountFailedAttemptsByUsernameAndIP(ctx interface{}, username interface{}, ipAddress interface{}, since interface{}) *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call {
	return &MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call{Call: _e.mock.On("CountFailedAttemptsByUsernameAndIP", ctx, username, ipAddress, since)}
}

func (_c *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call) Run(run func(ctx context.Context, username string, ipAddress string, since time.Time)) *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Time))
	})
	return _c
}

func (_c *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call) Return(_a0 int64, _a1 error) *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call) RunAndReturn(run func(context.Context, string, string, time.Time) (int64, error)) *MockLoginAttemptRepository_CountFailedAttemptsByUsernameAndIP_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, attempt
func (_m *MockLoginAttemptRepository) Create(ctx context.Context, attempt *domain.LoginAttempt) error {
	ret := _m.Called(ctx, attempt)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.LoginAttempt) error); ok {
		r0 = rf(ctx, attempt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockLoginAttemptRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockLoginAttemptRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - attempt *domain.LoginAttempt
func (_e *MockLoginAttemptRepository_Expecter) Create(ctx interface{}, attempt interface{}) *MockLoginAttemptRepository_Create_Call {
	return &MockLoginAttemptRepository_Create_Call{Call: _e.mock.On("Create", ctx, attempt)}
}

func (_c *MockLoginAttemptRepository_Create_Call) Run(run func(ctx context.Context, attempt *domain.LoginAttempt)) *MockLoginAttemptRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.LoginAttempt))
	})
	return _c
}

func (_c *MockLoginAttemptRepository_Create_Call) Return(_a0 error) *MockLoginAttemptRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLoginAttemptRepository_Create_Call) RunAndReturn(run func(context.Context, *domain.LoginAttempt) error) *MockLoginAttemptRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// CreateInTx provides a mock function with given fields: tx, attempt
func (_m *MockLoginAttemptRepository) CreateInTx(tx *gorm.DB, attempt *domain.LoginAttempt) error {
	ret := _m.Called(tx, attempt)

	if len(ret) == 0 {
		panic("no return value specified for CreateInTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*gorm.DB, *domain.LoginAttempt) error); ok {
		r0 = rf(tx, attempt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockLoginAttemptRepository_CreateInTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateInTx'
type MockLoginAttemptRepository_CreateInTx_Call struct {
	*mock.Call
}

// CreateInTx is a helper method to define mock.On call
//   - tx *gorm.DB
//   - attempt *domain.LoginAttempt
func (_e *MockLoginAttemptRepository_Expecter) CreateInTx(tx interface{}, attempt interface{}) *MockLoginAttemptRepository_CreateInTx_Call {
	return &MockLoginAttemptRepository_CreateInTx_Call{Call: _e.mock.On("CreateInTx", tx, attempt)}
}

func (_c *MockLoginAttemptRepository_CreateInTx_Call) Run(run func(tx *gorm.DB, attempt *domain.LoginAttempt)) *MockLoginAttemptRepository_CreateInTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gorm.DB), args[1].(*domain.LoginAttempt))
	})
	return _c
}

func (_c *MockLoginAttemptRepository_CreateInTx_Call) Return(_a0 error) *MockLoginAttemptRepository_CreateInTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLoginAttemptRepository_CreateInTx_Call) RunAndReturn(run func(*gorm.DB, *domain.LoginAttempt) error) *MockLoginAttemptRepository_CreateInTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLoginAttemptRepository creates a new instance of MockLoginAttemptRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLoginAttemptRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLoginAttemptRepository {
	mock := &MockLoginAttemptRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
