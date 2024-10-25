// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ignite/network/x/monitoringc/types"
)

// MonitoringcClient is an autogenerated mock type for the MonitoringcClient type
type MonitoringcClient struct {
	mock.Mock
}

type MonitoringcClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MonitoringcClient) EXPECT() *MonitoringcClient_Expecter {
	return &MonitoringcClient_Expecter{mock: &_m.Mock}
}

// GetLaunchIDFromChannelID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) GetLaunchIDFromChannelID(ctx context.Context, in *types.QueryGetLaunchIDFromChannelIDRequest, opts ...grpc.CallOption) (*types.QueryGetLaunchIDFromChannelIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetLaunchIDFromChannelID")
	}

	var r0 *types.QueryGetLaunchIDFromChannelIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetLaunchIDFromChannelIDRequest, ...grpc.CallOption) (*types.QueryGetLaunchIDFromChannelIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetLaunchIDFromChannelIDRequest, ...grpc.CallOption) *types.QueryGetLaunchIDFromChannelIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetLaunchIDFromChannelIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetLaunchIDFromChannelIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_GetLaunchIDFromChannelID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLaunchIDFromChannelID'
type MonitoringcClient_GetLaunchIDFromChannelID_Call struct {
	*mock.Call
}

// GetLaunchIDFromChannelID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetLaunchIDFromChannelIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) GetLaunchIDFromChannelID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_GetLaunchIDFromChannelID_Call {
	return &MonitoringcClient_GetLaunchIDFromChannelID_Call{Call: _e.mock.On("GetLaunchIDFromChannelID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_GetLaunchIDFromChannelID_Call) Run(run func(ctx context.Context, in *types.QueryGetLaunchIDFromChannelIDRequest, opts ...grpc.CallOption)) *MonitoringcClient_GetLaunchIDFromChannelID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetLaunchIDFromChannelIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_GetLaunchIDFromChannelID_Call) Return(_a0 *types.QueryGetLaunchIDFromChannelIDResponse, _a1 error) *MonitoringcClient_GetLaunchIDFromChannelID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_GetLaunchIDFromChannelID_Call) RunAndReturn(run func(context.Context, *types.QueryGetLaunchIDFromChannelIDRequest, ...grpc.CallOption) (*types.QueryGetLaunchIDFromChannelIDResponse, error)) *MonitoringcClient_GetLaunchIDFromChannelID_Call {
	_c.Call.Return(run)
	return _c
}

// GetMonitoringHistory provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) GetMonitoringHistory(ctx context.Context, in *types.QueryGetMonitoringHistoryRequest, opts ...grpc.CallOption) (*types.QueryGetMonitoringHistoryResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetMonitoringHistory")
	}

	var r0 *types.QueryGetMonitoringHistoryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetMonitoringHistoryRequest, ...grpc.CallOption) (*types.QueryGetMonitoringHistoryResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetMonitoringHistoryRequest, ...grpc.CallOption) *types.QueryGetMonitoringHistoryResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetMonitoringHistoryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetMonitoringHistoryRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_GetMonitoringHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMonitoringHistory'
type MonitoringcClient_GetMonitoringHistory_Call struct {
	*mock.Call
}

// GetMonitoringHistory is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetMonitoringHistoryRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) GetMonitoringHistory(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_GetMonitoringHistory_Call {
	return &MonitoringcClient_GetMonitoringHistory_Call{Call: _e.mock.On("GetMonitoringHistory",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_GetMonitoringHistory_Call) Run(run func(ctx context.Context, in *types.QueryGetMonitoringHistoryRequest, opts ...grpc.CallOption)) *MonitoringcClient_GetMonitoringHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetMonitoringHistoryRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_GetMonitoringHistory_Call) Return(_a0 *types.QueryGetMonitoringHistoryResponse, _a1 error) *MonitoringcClient_GetMonitoringHistory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_GetMonitoringHistory_Call) RunAndReturn(run func(context.Context, *types.QueryGetMonitoringHistoryRequest, ...grpc.CallOption) (*types.QueryGetMonitoringHistoryResponse, error)) *MonitoringcClient_GetMonitoringHistory_Call {
	_c.Call.Return(run)
	return _c
}

// GetProviderClientID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) GetProviderClientID(ctx context.Context, in *types.QueryGetProviderClientIDRequest, opts ...grpc.CallOption) (*types.QueryGetProviderClientIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetProviderClientID")
	}

	var r0 *types.QueryGetProviderClientIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetProviderClientIDRequest, ...grpc.CallOption) (*types.QueryGetProviderClientIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetProviderClientIDRequest, ...grpc.CallOption) *types.QueryGetProviderClientIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetProviderClientIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetProviderClientIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_GetProviderClientID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProviderClientID'
type MonitoringcClient_GetProviderClientID_Call struct {
	*mock.Call
}

// GetProviderClientID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetProviderClientIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) GetProviderClientID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_GetProviderClientID_Call {
	return &MonitoringcClient_GetProviderClientID_Call{Call: _e.mock.On("GetProviderClientID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_GetProviderClientID_Call) Run(run func(ctx context.Context, in *types.QueryGetProviderClientIDRequest, opts ...grpc.CallOption)) *MonitoringcClient_GetProviderClientID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetProviderClientIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_GetProviderClientID_Call) Return(_a0 *types.QueryGetProviderClientIDResponse, _a1 error) *MonitoringcClient_GetProviderClientID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_GetProviderClientID_Call) RunAndReturn(run func(context.Context, *types.QueryGetProviderClientIDRequest, ...grpc.CallOption) (*types.QueryGetProviderClientIDResponse, error)) *MonitoringcClient_GetProviderClientID_Call {
	_c.Call.Return(run)
	return _c
}

// GetVerifiedClientID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) GetVerifiedClientID(ctx context.Context, in *types.QueryGetVerifiedClientIDRequest, opts ...grpc.CallOption) (*types.QueryGetVerifiedClientIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetVerifiedClientID")
	}

	var r0 *types.QueryGetVerifiedClientIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetVerifiedClientIDRequest, ...grpc.CallOption) (*types.QueryGetVerifiedClientIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetVerifiedClientIDRequest, ...grpc.CallOption) *types.QueryGetVerifiedClientIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetVerifiedClientIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetVerifiedClientIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_GetVerifiedClientID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVerifiedClientID'
type MonitoringcClient_GetVerifiedClientID_Call struct {
	*mock.Call
}

// GetVerifiedClientID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetVerifiedClientIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) GetVerifiedClientID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_GetVerifiedClientID_Call {
	return &MonitoringcClient_GetVerifiedClientID_Call{Call: _e.mock.On("GetVerifiedClientID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_GetVerifiedClientID_Call) Run(run func(ctx context.Context, in *types.QueryGetVerifiedClientIDRequest, opts ...grpc.CallOption)) *MonitoringcClient_GetVerifiedClientID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetVerifiedClientIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_GetVerifiedClientID_Call) Return(_a0 *types.QueryGetVerifiedClientIDResponse, _a1 error) *MonitoringcClient_GetVerifiedClientID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_GetVerifiedClientID_Call) RunAndReturn(run func(context.Context, *types.QueryGetVerifiedClientIDRequest, ...grpc.CallOption) (*types.QueryGetVerifiedClientIDResponse, error)) *MonitoringcClient_GetVerifiedClientID_Call {
	_c.Call.Return(run)
	return _c
}

// ListLaunchIDFromChannelID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) ListLaunchIDFromChannelID(ctx context.Context, in *types.QueryAllLaunchIDFromChannelIDRequest, opts ...grpc.CallOption) (*types.QueryAllLaunchIDFromChannelIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListLaunchIDFromChannelID")
	}

	var r0 *types.QueryAllLaunchIDFromChannelIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllLaunchIDFromChannelIDRequest, ...grpc.CallOption) (*types.QueryAllLaunchIDFromChannelIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllLaunchIDFromChannelIDRequest, ...grpc.CallOption) *types.QueryAllLaunchIDFromChannelIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllLaunchIDFromChannelIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllLaunchIDFromChannelIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_ListLaunchIDFromChannelID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListLaunchIDFromChannelID'
type MonitoringcClient_ListLaunchIDFromChannelID_Call struct {
	*mock.Call
}

// ListLaunchIDFromChannelID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryAllLaunchIDFromChannelIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) ListLaunchIDFromChannelID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_ListLaunchIDFromChannelID_Call {
	return &MonitoringcClient_ListLaunchIDFromChannelID_Call{Call: _e.mock.On("ListLaunchIDFromChannelID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_ListLaunchIDFromChannelID_Call) Run(run func(ctx context.Context, in *types.QueryAllLaunchIDFromChannelIDRequest, opts ...grpc.CallOption)) *MonitoringcClient_ListLaunchIDFromChannelID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryAllLaunchIDFromChannelIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_ListLaunchIDFromChannelID_Call) Return(_a0 *types.QueryAllLaunchIDFromChannelIDResponse, _a1 error) *MonitoringcClient_ListLaunchIDFromChannelID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_ListLaunchIDFromChannelID_Call) RunAndReturn(run func(context.Context, *types.QueryAllLaunchIDFromChannelIDRequest, ...grpc.CallOption) (*types.QueryAllLaunchIDFromChannelIDResponse, error)) *MonitoringcClient_ListLaunchIDFromChannelID_Call {
	_c.Call.Return(run)
	return _c
}

// ListProviderClientID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) ListProviderClientID(ctx context.Context, in *types.QueryAllProviderClientIDRequest, opts ...grpc.CallOption) (*types.QueryAllProviderClientIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListProviderClientID")
	}

	var r0 *types.QueryAllProviderClientIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllProviderClientIDRequest, ...grpc.CallOption) (*types.QueryAllProviderClientIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllProviderClientIDRequest, ...grpc.CallOption) *types.QueryAllProviderClientIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllProviderClientIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllProviderClientIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_ListProviderClientID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListProviderClientID'
type MonitoringcClient_ListProviderClientID_Call struct {
	*mock.Call
}

// ListProviderClientID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryAllProviderClientIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) ListProviderClientID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_ListProviderClientID_Call {
	return &MonitoringcClient_ListProviderClientID_Call{Call: _e.mock.On("ListProviderClientID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_ListProviderClientID_Call) Run(run func(ctx context.Context, in *types.QueryAllProviderClientIDRequest, opts ...grpc.CallOption)) *MonitoringcClient_ListProviderClientID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryAllProviderClientIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_ListProviderClientID_Call) Return(_a0 *types.QueryAllProviderClientIDResponse, _a1 error) *MonitoringcClient_ListProviderClientID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_ListProviderClientID_Call) RunAndReturn(run func(context.Context, *types.QueryAllProviderClientIDRequest, ...grpc.CallOption) (*types.QueryAllProviderClientIDResponse, error)) *MonitoringcClient_ListProviderClientID_Call {
	_c.Call.Return(run)
	return _c
}

// Params provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringcClient) Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Params")
	}

	var r0 *types.QueryParamsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) (*types.QueryParamsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) *types.QueryParamsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryParamsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringcClient_Params_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Params'
type MonitoringcClient_Params_Call struct {
	*mock.Call
}

// Params is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryParamsRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringcClient_Expecter) Params(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringcClient_Params_Call {
	return &MonitoringcClient_Params_Call{Call: _e.mock.On("Params",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringcClient_Params_Call) Run(run func(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption)) *MonitoringcClient_Params_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryParamsRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringcClient_Params_Call) Return(_a0 *types.QueryParamsResponse, _a1 error) *MonitoringcClient_Params_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringcClient_Params_Call) RunAndReturn(run func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) (*types.QueryParamsResponse, error)) *MonitoringcClient_Params_Call {
	_c.Call.Return(run)
	return _c
}

// NewMonitoringcClient creates a new instance of MonitoringcClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMonitoringcClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MonitoringcClient {
	mock := &MonitoringcClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
