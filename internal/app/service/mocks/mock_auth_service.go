// Code generated by mockery v2.53.4. DO NOT EDIT.

package service

import (
	context "context"
	domain "beerdosan-backend/internal/app/domain"

	mock "github.com/stretchr/testify/mock"

	service "beerdosan-backend/internal/app/service"
)

// MockAuthService is an autogenerated mock type for the AuthService type
type MockAuthService struct {
	mock.Mock
}

type MockAuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthService) EXPECT() *MockAuthService_Expecter {
	return &MockAuthService_Expecter{mock: &_m.Mock}
}

// CheckRateLimit provides a mock function with given fields: ctx, username, ipAddress
func (_m *MockAuthService) CheckRateLimit(ctx context.Context, username string, ipAddress string) error {
	ret := _m.Called(ctx, username, ipAddress)

	if len(ret) == 0 {
		panic("no return value specified for CheckRateLimit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, username, ipAddress)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_CheckRateLimit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckRateLimit'
type MockAuthService_CheckRateLimit_Call struct {
	*mock.Call
}

// CheckRateLimit is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - ipAddress string
func (_e *MockAuthService_Expecter) CheckRateLimit(ctx interface{}, username interface{}, ipAddress interface{}) *MockAuthService_CheckRateLimit_Call {
	return &MockAuthService_CheckRateLimit_Call{Call: _e.mock.On("CheckRateLimit", ctx, username, ipAddress)}
}

func (_c *MockAuthService_CheckRateLimit_Call) Run(run func(ctx context.Context, username string, ipAddress string)) *MockAuthService_CheckRateLimit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthService_CheckRateLimit_Call) Return(_a0 error) *MockAuthService_CheckRateLimit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_CheckRateLimit_Call) RunAndReturn(run func(context.Context, string, string) error) *MockAuthService_CheckRateLimit_Call {
	_c.Call.Return(run)
	return _c
}

// CleanupExpiredSessions provides a mock function with given fields: ctx
func (_m *MockAuthService) CleanupExpiredSessions(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CleanupExpiredSessions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_CleanupExpiredSessions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CleanupExpiredSessions'
type MockAuthService_CleanupExpiredSessions_Call struct {
	*mock.Call
}

// CleanupExpiredSessions is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockAuthService_Expecter) CleanupExpiredSessions(ctx interface{}) *MockAuthService_CleanupExpiredSessions_Call {
	return &MockAuthService_CleanupExpiredSessions_Call{Call: _e.mock.On("CleanupExpiredSessions", ctx)}
}

func (_c *MockAuthService_CleanupExpiredSessions_Call) Run(run func(ctx context.Context)) *MockAuthService_CleanupExpiredSessions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockAuthService_CleanupExpiredSessions_Call) Return(_a0 error) *MockAuthService_CleanupExpiredSessions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_CleanupExpiredSessions_Call) RunAndReturn(run func(context.Context) error) *MockAuthService_CleanupExpiredSessions_Call {
	_c.Call.Return(run)
	return _c
}

// CreateSession provides a mock function with given fields: ctx, userID, deviceInfo, ipAddress
func (_m *MockAuthService) CreateSession(ctx context.Context, userID domain.UserID, deviceInfo string, ipAddress string) (*domain.Session, error) {
	ret := _m.Called(ctx, userID, deviceInfo, ipAddress)

	if len(ret) == 0 {
		panic("no return value specified for CreateSession")
	}

	var r0 *domain.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserID, string, string) (*domain.Session, error)); ok {
		return rf(ctx, userID, deviceInfo, ipAddress)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserID, string, string) *domain.Session); ok {
		r0 = rf(ctx, userID, deviceInfo, ipAddress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.UserID, string, string) error); ok {
		r1 = rf(ctx, userID, deviceInfo, ipAddress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthService_CreateSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateSession'
type MockAuthService_CreateSession_Call struct {
	*mock.Call
}

// CreateSession is a helper method to define mock.On call
//   - ctx context.Context
//   - userID domain.UserID
//   - deviceInfo string
//   - ipAddress string
func (_e *MockAuthService_Expecter) CreateSession(ctx interface{}, userID interface{}, deviceInfo interface{}, ipAddress interface{}) *MockAuthService_CreateSession_Call {
	return &MockAuthService_CreateSession_Call{Call: _e.mock.On("CreateSession", ctx, userID, deviceInfo, ipAddress)}
}

func (_c *MockAuthService_CreateSession_Call) Run(run func(ctx context.Context, userID domain.UserID, deviceInfo string, ipAddress string)) *MockAuthService_CreateSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.UserID), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *MockAuthService_CreateSession_Call) Return(_a0 *domain.Session, _a1 error) *MockAuthService_CreateSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthService_CreateSession_Call) RunAndReturn(run func(context.Context, domain.UserID, string, string) (*domain.Session, error)) *MockAuthService_CreateSession_Call {
	_c.Call.Return(run)
	return _c
}

// InvalidateAllUserSessions provides a mock function with given fields: ctx, userID, excludeSessionID
func (_m *MockAuthService) InvalidateAllUserSessions(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID) error {
	ret := _m.Called(ctx, userID, excludeSessionID)

	if len(ret) == 0 {
		panic("no return value specified for InvalidateAllUserSessions")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.UserID, domain.SessionID) error); ok {
		r0 = rf(ctx, userID, excludeSessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_InvalidateAllUserSessions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InvalidateAllUserSessions'
type MockAuthService_InvalidateAllUserSessions_Call struct {
	*mock.Call
}

// InvalidateAllUserSessions is a helper method to define mock.On call
//   - ctx context.Context
//   - userID domain.UserID
//   - excludeSessionID domain.SessionID
func (_e *MockAuthService_Expecter) InvalidateAllUserSessions(ctx interface{}, userID interface{}, excludeSessionID interface{}) *MockAuthService_InvalidateAllUserSessions_Call {
	return &MockAuthService_InvalidateAllUserSessions_Call{Call: _e.mock.On("InvalidateAllUserSessions", ctx, userID, excludeSessionID)}
}

func (_c *MockAuthService_InvalidateAllUserSessions_Call) Run(run func(ctx context.Context, userID domain.UserID, excludeSessionID domain.SessionID)) *MockAuthService_InvalidateAllUserSessions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.UserID), args[2].(domain.SessionID))
	})
	return _c
}

func (_c *MockAuthService_InvalidateAllUserSessions_Call) Return(_a0 error) *MockAuthService_InvalidateAllUserSessions_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_InvalidateAllUserSessions_Call) RunAndReturn(run func(context.Context, domain.UserID, domain.SessionID) error) *MockAuthService_InvalidateAllUserSessions_Call {
	_c.Call.Return(run)
	return _c
}

// InvalidateSession provides a mock function with given fields: ctx, sessionID
func (_m *MockAuthService) InvalidateSession(ctx context.Context, sessionID domain.SessionID) error {
	ret := _m.Called(ctx, sessionID)

	if len(ret) == 0 {
		panic("no return value specified for InvalidateSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.SessionID) error); ok {
		r0 = rf(ctx, sessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_InvalidateSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InvalidateSession'
type MockAuthService_InvalidateSession_Call struct {
	*mock.Call
}

// InvalidateSession is a helper method to define mock.On call
//   - ctx context.Context
//   - sessionID domain.SessionID
func (_e *MockAuthService_Expecter) InvalidateSession(ctx interface{}, sessionID interface{}) *MockAuthService_InvalidateSession_Call {
	return &MockAuthService_InvalidateSession_Call{Call: _e.mock.On("InvalidateSession", ctx, sessionID)}
}

func (_c *MockAuthService_InvalidateSession_Call) Run(run func(ctx context.Context, sessionID domain.SessionID)) *MockAuthService_InvalidateSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.SessionID))
	})
	return _c
}

func (_c *MockAuthService_InvalidateSession_Call) Return(_a0 error) *MockAuthService_InvalidateSession_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_InvalidateSession_Call) RunAndReturn(run func(context.Context, domain.SessionID) error) *MockAuthService_InvalidateSession_Call {
	_c.Call.Return(run)
	return _c
}

// RecordLoginAttempt provides a mock function with given fields: ctx, username, ipAddress, success, failureReason
func (_m *MockAuthService) RecordLoginAttempt(ctx context.Context, username string, ipAddress string, success bool, failureReason string) error {
	ret := _m.Called(ctx, username, ipAddress, success, failureReason)

	if len(ret) == 0 {
		panic("no return value specified for RecordLoginAttempt")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool, string) error); ok {
		r0 = rf(ctx, username, ipAddress, success, failureReason)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_RecordLoginAttempt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecordLoginAttempt'
type MockAuthService_RecordLoginAttempt_Call struct {
	*mock.Call
}

// RecordLoginAttempt is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - ipAddress string
//   - success bool
//   - failureReason string
func (_e *MockAuthService_Expecter) RecordLoginAttempt(ctx interface{}, username interface{}, ipAddress interface{}, success interface{}, failureReason interface{}) *MockAuthService_RecordLoginAttempt_Call {
	return &MockAuthService_RecordLoginAttempt_Call{Call: _e.mock.On("RecordLoginAttempt", ctx, username, ipAddress, success, failureReason)}
}

func (_c *MockAuthService_RecordLoginAttempt_Call) Run(run func(ctx context.Context, username string, ipAddress string, success bool, failureReason string)) *MockAuthService_RecordLoginAttempt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(bool), args[4].(string))
	})
	return _c
}

func (_c *MockAuthService_RecordLoginAttempt_Call) Return(_a0 error) *MockAuthService_RecordLoginAttempt_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_RecordLoginAttempt_Call) RunAndReturn(run func(context.Context, string, string, bool, string) error) *MockAuthService_RecordLoginAttempt_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSessionActivity provides a mock function with given fields: ctx, sessionID
func (_m *MockAuthService) UpdateSessionActivity(ctx context.Context, sessionID domain.SessionID) error {
	ret := _m.Called(ctx, sessionID)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSessionActivity")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.SessionID) error); ok {
		r0 = rf(ctx, sessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthService_UpdateSessionActivity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSessionActivity'
type MockAuthService_UpdateSessionActivity_Call struct {
	*mock.Call
}

// UpdateSessionActivity is a helper method to define mock.On call
//   - ctx context.Context
//   - sessionID domain.SessionID
func (_e *MockAuthService_Expecter) UpdateSessionActivity(ctx interface{}, sessionID interface{}) *MockAuthService_UpdateSessionActivity_Call {
	return &MockAuthService_UpdateSessionActivity_Call{Call: _e.mock.On("UpdateSessionActivity", ctx, sessionID)}
}

func (_c *MockAuthService_UpdateSessionActivity_Call) Run(run func(ctx context.Context, sessionID domain.SessionID)) *MockAuthService_UpdateSessionActivity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.SessionID))
	})
	return _c
}

func (_c *MockAuthService_UpdateSessionActivity_Call) Return(_a0 error) *MockAuthService_UpdateSessionActivity_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthService_UpdateSessionActivity_Call) RunAndReturn(run func(context.Context, domain.SessionID) error) *MockAuthService_UpdateSessionActivity_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateCredentials provides a mock function with given fields: ctx, username, password
func (_m *MockAuthService) ValidateCredentials(ctx context.Context, username string, password string) (*domain.User, error) {
	ret := _m.Called(ctx, username, password)

	if len(ret) == 0 {
		panic("no return value specified for ValidateCredentials")
	}

	var r0 *domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*domain.User, error)); ok {
		return rf(ctx, username, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *domain.User); ok {
		r0 = rf(ctx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthService_ValidateCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateCredentials'
type MockAuthService_ValidateCredentials_Call struct {
	*mock.Call
}

// ValidateCredentials is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - password string
func (_e *MockAuthService_Expecter) ValidateCredentials(ctx interface{}, username interface{}, password interface{}) *MockAuthService_ValidateCredentials_Call {
	return &MockAuthService_ValidateCredentials_Call{Call: _e.mock.On("ValidateCredentials", ctx, username, password)}
}

func (_c *MockAuthService_ValidateCredentials_Call) Run(run func(ctx context.Context, username string, password string)) *MockAuthService_ValidateCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockAuthService_ValidateCredentials_Call) Return(_a0 *domain.User, _a1 error) *MockAuthService_ValidateCredentials_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthService_ValidateCredentials_Call) RunAndReturn(run func(context.Context, string, string) (*domain.User, error)) *MockAuthService_ValidateCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateSession provides a mock function with given fields: ctx, sessionID
func (_m *MockAuthService) ValidateSession(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	ret := _m.Called(ctx, sessionID)

	if len(ret) == 0 {
		panic("no return value specified for ValidateSession")
	}

	var r0 *domain.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.SessionID) (*domain.Session, error)); ok {
		return rf(ctx, sessionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.SessionID) *domain.Session); ok {
		r0 = rf(ctx, sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.SessionID) error); ok {
		r1 = rf(ctx, sessionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthService_ValidateSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateSession'
type MockAuthService_ValidateSession_Call struct {
	*mock.Call
}

// ValidateSession is a helper method to define mock.On call
//   - ctx context.Context
//   - sessionID domain.SessionID
func (_e *MockAuthService_Expecter) ValidateSession(ctx interface{}, sessionID interface{}) *MockAuthService_ValidateSession_Call {
	return &MockAuthService_ValidateSession_Call{Call: _e.mock.On("ValidateSession", ctx, sessionID)}
}

func (_c *MockAuthService_ValidateSession_Call) Run(run func(ctx context.Context, sessionID domain.SessionID)) *MockAuthService_ValidateSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.SessionID))
	})
	return _c
}

func (_c *MockAuthService_ValidateSession_Call) Return(_a0 *domain.Session, _a1 error) *MockAuthService_ValidateSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthService_ValidateSession_Call) RunAndReturn(run func(context.Context, domain.SessionID) (*domain.Session, error)) *MockAuthService_ValidateSession_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateToken provides a mock function with given fields: ctx, token
func (_m *MockAuthService) ValidateToken(ctx context.Context, token string) (*service.AuthClaims, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *service.AuthClaims
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*service.AuthClaims, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *service.AuthClaims); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.AuthClaims)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthService_ValidateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateToken'
type MockAuthService_ValidateToken_Call struct {
	*mock.Call
}

// ValidateToken is a helper method to define mock.On call
//   - ctx context.Context
//   - token string
func (_e *MockAuthService_Expecter) ValidateToken(ctx interface{}, token interface{}) *MockAuthService_ValidateToken_Call {
	return &MockAuthService_ValidateToken_Call{Call: _e.mock.On("ValidateToken", ctx, token)}
}

func (_c *MockAuthService_ValidateToken_Call) Run(run func(ctx context.Context, token string)) *MockAuthService_ValidateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockAuthService_ValidateToken_Call) Return(_a0 *service.AuthClaims, _a1 error) *MockAuthService_ValidateToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthService_ValidateToken_Call) RunAndReturn(run func(context.Context, string) (*service.AuthClaims, error)) *MockAuthService_ValidateToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAuthService creates a new instance of MockAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthService {
	mock := &MockAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
