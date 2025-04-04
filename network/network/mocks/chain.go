// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Chain is an autogenerated mock type for the Chain type
type Chain struct {
	mock.Mock
}

type Chain_Expecter struct {
	mock *mock.Mock
}

func (_m *Chain) EXPECT() *Chain_Expecter {
	return &Chain_Expecter{mock: &_m.Mock}
}

// AppTOMLPath provides a mock function with given fields:
func (_m *Chain) AppTOMLPath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AppTOMLPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_AppTOMLPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppTOMLPath'
type Chain_AppTOMLPath_Call struct {
	*mock.Call
}

// AppTOMLPath is a helper method to define mock.On call
func (_e *Chain_Expecter) AppTOMLPath() *Chain_AppTOMLPath_Call {
	return &Chain_AppTOMLPath_Call{Call: _e.mock.On("AppTOMLPath")}
}

func (_c *Chain_AppTOMLPath_Call) Run(run func()) *Chain_AppTOMLPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_AppTOMLPath_Call) Return(_a0 string, _a1 error) *Chain_AppTOMLPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_AppTOMLPath_Call) RunAndReturn(run func() (string, error)) *Chain_AppTOMLPath_Call {
	_c.Call.Return(run)
	return _c
}

// CacheBinary provides a mock function with given fields: launchID
func (_m *Chain) CacheBinary(launchID uint64) error {
	ret := _m.Called(launchID)

	if len(ret) == 0 {
		panic("no return value specified for CacheBinary")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(launchID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Chain_CacheBinary_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CacheBinary'
type Chain_CacheBinary_Call struct {
	*mock.Call
}

// CacheBinary is a helper method to define mock.On call
//   - launchID uint64
func (_e *Chain_Expecter) CacheBinary(launchID interface{}) *Chain_CacheBinary_Call {
	return &Chain_CacheBinary_Call{Call: _e.mock.On("CacheBinary", launchID)}
}

func (_c *Chain_CacheBinary_Call) Run(run func(launchID uint64)) *Chain_CacheBinary_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint64))
	})
	return _c
}

func (_c *Chain_CacheBinary_Call) Return(_a0 error) *Chain_CacheBinary_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Chain_CacheBinary_Call) RunAndReturn(run func(uint64) error) *Chain_CacheBinary_Call {
	_c.Call.Return(run)
	return _c
}

// ChainID provides a mock function with given fields:
func (_m *Chain) ChainID() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ChainID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_ChainID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChainID'
type Chain_ChainID_Call struct {
	*mock.Call
}

// ChainID is a helper method to define mock.On call
func (_e *Chain_Expecter) ChainID() *Chain_ChainID_Call {
	return &Chain_ChainID_Call{Call: _e.mock.On("ChainID")}
}

func (_c *Chain_ChainID_Call) Run(run func()) *Chain_ChainID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_ChainID_Call) Return(_a0 string, _a1 error) *Chain_ChainID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_ChainID_Call) RunAndReturn(run func() (string, error)) *Chain_ChainID_Call {
	_c.Call.Return(run)
	return _c
}

// ConfigTOMLPath provides a mock function with given fields:
func (_m *Chain) ConfigTOMLPath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ConfigTOMLPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_ConfigTOMLPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfigTOMLPath'
type Chain_ConfigTOMLPath_Call struct {
	*mock.Call
}

// ConfigTOMLPath is a helper method to define mock.On call
func (_e *Chain_Expecter) ConfigTOMLPath() *Chain_ConfigTOMLPath_Call {
	return &Chain_ConfigTOMLPath_Call{Call: _e.mock.On("ConfigTOMLPath")}
}

func (_c *Chain_ConfigTOMLPath_Call) Run(run func()) *Chain_ConfigTOMLPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_ConfigTOMLPath_Call) Return(_a0 string, _a1 error) *Chain_ConfigTOMLPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_ConfigTOMLPath_Call) RunAndReturn(run func() (string, error)) *Chain_ConfigTOMLPath_Call {
	_c.Call.Return(run)
	return _c
}

// DefaultGentxPath provides a mock function with given fields:
func (_m *Chain) DefaultGentxPath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for DefaultGentxPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_DefaultGentxPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DefaultGentxPath'
type Chain_DefaultGentxPath_Call struct {
	*mock.Call
}

// DefaultGentxPath is a helper method to define mock.On call
func (_e *Chain_Expecter) DefaultGentxPath() *Chain_DefaultGentxPath_Call {
	return &Chain_DefaultGentxPath_Call{Call: _e.mock.On("DefaultGentxPath")}
}

func (_c *Chain_DefaultGentxPath_Call) Run(run func()) *Chain_DefaultGentxPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_DefaultGentxPath_Call) Return(_a0 string, _a1 error) *Chain_DefaultGentxPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_DefaultGentxPath_Call) RunAndReturn(run func() (string, error)) *Chain_DefaultGentxPath_Call {
	_c.Call.Return(run)
	return _c
}

// GenesisPath provides a mock function with given fields:
func (_m *Chain) GenesisPath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GenesisPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_GenesisPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenesisPath'
type Chain_GenesisPath_Call struct {
	*mock.Call
}

// GenesisPath is a helper method to define mock.On call
func (_e *Chain_Expecter) GenesisPath() *Chain_GenesisPath_Call {
	return &Chain_GenesisPath_Call{Call: _e.mock.On("GenesisPath")}
}

func (_c *Chain_GenesisPath_Call) Run(run func()) *Chain_GenesisPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_GenesisPath_Call) Return(_a0 string, _a1 error) *Chain_GenesisPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_GenesisPath_Call) RunAndReturn(run func() (string, error)) *Chain_GenesisPath_Call {
	_c.Call.Return(run)
	return _c
}

// GentxsPath provides a mock function with given fields:
func (_m *Chain) GentxsPath() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GentxsPath")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_GentxsPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GentxsPath'
type Chain_GentxsPath_Call struct {
	*mock.Call
}

// GentxsPath is a helper method to define mock.On call
func (_e *Chain_Expecter) GentxsPath() *Chain_GentxsPath_Call {
	return &Chain_GentxsPath_Call{Call: _e.mock.On("GentxsPath")}
}

func (_c *Chain_GentxsPath_Call) Run(run func()) *Chain_GentxsPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_GentxsPath_Call) Return(_a0 string, _a1 error) *Chain_GentxsPath_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_GentxsPath_Call) RunAndReturn(run func() (string, error)) *Chain_GentxsPath_Call {
	_c.Call.Return(run)
	return _c
}

// ID provides a mock function with given fields:
func (_m *Chain) ID() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_ID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ID'
type Chain_ID_Call struct {
	*mock.Call
}

// ID is a helper method to define mock.On call
func (_e *Chain_Expecter) ID() *Chain_ID_Call {
	return &Chain_ID_Call{Call: _e.mock.On("ID")}
}

func (_c *Chain_ID_Call) Run(run func()) *Chain_ID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_ID_Call) Return(_a0 string, _a1 error) *Chain_ID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_ID_Call) RunAndReturn(run func() (string, error)) *Chain_ID_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *Chain) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Chain_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type Chain_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *Chain_Expecter) Name() *Chain_Name_Call {
	return &Chain_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *Chain_Name_Call) Run(run func()) *Chain_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_Name_Call) Return(_a0 string) *Chain_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Chain_Name_Call) RunAndReturn(run func() string) *Chain_Name_Call {
	_c.Call.Return(run)
	return _c
}

// NodeID provides a mock function with given fields: ctx
func (_m *Chain) NodeID(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for NodeID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Chain_NodeID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NodeID'
type Chain_NodeID_Call struct {
	*mock.Call
}

// NodeID is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Chain_Expecter) NodeID(ctx interface{}) *Chain_NodeID_Call {
	return &Chain_NodeID_Call{Call: _e.mock.On("NodeID", ctx)}
}

func (_c *Chain_NodeID_Call) Run(run func(ctx context.Context)) *Chain_NodeID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Chain_NodeID_Call) Return(_a0 string, _a1 error) *Chain_NodeID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Chain_NodeID_Call) RunAndReturn(run func(context.Context) (string, error)) *Chain_NodeID_Call {
	_c.Call.Return(run)
	return _c
}

// SourceHash provides a mock function with given fields:
func (_m *Chain) SourceHash() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SourceHash")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Chain_SourceHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SourceHash'
type Chain_SourceHash_Call struct {
	*mock.Call
}

// SourceHash is a helper method to define mock.On call
func (_e *Chain_Expecter) SourceHash() *Chain_SourceHash_Call {
	return &Chain_SourceHash_Call{Call: _e.mock.On("SourceHash")}
}

func (_c *Chain_SourceHash_Call) Run(run func()) *Chain_SourceHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_SourceHash_Call) Return(_a0 string) *Chain_SourceHash_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Chain_SourceHash_Call) RunAndReturn(run func() string) *Chain_SourceHash_Call {
	_c.Call.Return(run)
	return _c
}

// SourceURL provides a mock function with given fields:
func (_m *Chain) SourceURL() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SourceURL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Chain_SourceURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SourceURL'
type Chain_SourceURL_Call struct {
	*mock.Call
}

// SourceURL is a helper method to define mock.On call
func (_e *Chain_Expecter) SourceURL() *Chain_SourceURL_Call {
	return &Chain_SourceURL_Call{Call: _e.mock.On("SourceURL")}
}

func (_c *Chain_SourceURL_Call) Run(run func()) *Chain_SourceURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Chain_SourceURL_Call) Return(_a0 string) *Chain_SourceURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Chain_SourceURL_Call) RunAndReturn(run func() string) *Chain_SourceURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewChain creates a new instance of Chain. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChain(t interface {
	mock.TestingT
	Cleanup(func())
}) *Chain {
	mock := &Chain{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
