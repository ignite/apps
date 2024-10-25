// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ignite/network/x/monitoringp/types"
)

// MonitoringpClient is an autogenerated mock type for the MonitoringpClient type
type MonitoringpClient struct {
	mock.Mock
}

type MonitoringpClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MonitoringpClient) EXPECT() *MonitoringpClient_Expecter {
	return &MonitoringpClient_Expecter{mock: &_m.Mock}
}

// GetConnectionChannelID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringpClient) GetConnectionChannelID(ctx context.Context, in *types.QueryGetConnectionChannelIDRequest, opts ...grpc.CallOption) (*types.QueryGetConnectionChannelIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetConnectionChannelID")
	}

	var r0 *types.QueryGetConnectionChannelIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetConnectionChannelIDRequest, ...grpc.CallOption) (*types.QueryGetConnectionChannelIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetConnectionChannelIDRequest, ...grpc.CallOption) *types.QueryGetConnectionChannelIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetConnectionChannelIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetConnectionChannelIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringpClient_GetConnectionChannelID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnectionChannelID'
type MonitoringpClient_GetConnectionChannelID_Call struct {
	*mock.Call
}

// GetConnectionChannelID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetConnectionChannelIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringpClient_Expecter) GetConnectionChannelID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringpClient_GetConnectionChannelID_Call {
	return &MonitoringpClient_GetConnectionChannelID_Call{Call: _e.mock.On("GetConnectionChannelID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringpClient_GetConnectionChannelID_Call) Run(run func(ctx context.Context, in *types.QueryGetConnectionChannelIDRequest, opts ...grpc.CallOption)) *MonitoringpClient_GetConnectionChannelID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetConnectionChannelIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringpClient_GetConnectionChannelID_Call) Return(_a0 *types.QueryGetConnectionChannelIDResponse, _a1 error) *MonitoringpClient_GetConnectionChannelID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringpClient_GetConnectionChannelID_Call) RunAndReturn(run func(context.Context, *types.QueryGetConnectionChannelIDRequest, ...grpc.CallOption) (*types.QueryGetConnectionChannelIDResponse, error)) *MonitoringpClient_GetConnectionChannelID_Call {
	_c.Call.Return(run)
	return _c
}

// GetConsumerClientID provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringpClient) GetConsumerClientID(ctx context.Context, in *types.QueryGetConsumerClientIDRequest, opts ...grpc.CallOption) (*types.QueryGetConsumerClientIDResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetConsumerClientID")
	}

	var r0 *types.QueryGetConsumerClientIDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetConsumerClientIDRequest, ...grpc.CallOption) (*types.QueryGetConsumerClientIDResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetConsumerClientIDRequest, ...grpc.CallOption) *types.QueryGetConsumerClientIDResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetConsumerClientIDResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetConsumerClientIDRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringpClient_GetConsumerClientID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConsumerClientID'
type MonitoringpClient_GetConsumerClientID_Call struct {
	*mock.Call
}

// GetConsumerClientID is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetConsumerClientIDRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringpClient_Expecter) GetConsumerClientID(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringpClient_GetConsumerClientID_Call {
	return &MonitoringpClient_GetConsumerClientID_Call{Call: _e.mock.On("GetConsumerClientID",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringpClient_GetConsumerClientID_Call) Run(run func(ctx context.Context, in *types.QueryGetConsumerClientIDRequest, opts ...grpc.CallOption)) *MonitoringpClient_GetConsumerClientID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetConsumerClientIDRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringpClient_GetConsumerClientID_Call) Return(_a0 *types.QueryGetConsumerClientIDResponse, _a1 error) *MonitoringpClient_GetConsumerClientID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringpClient_GetConsumerClientID_Call) RunAndReturn(run func(context.Context, *types.QueryGetConsumerClientIDRequest, ...grpc.CallOption) (*types.QueryGetConsumerClientIDResponse, error)) *MonitoringpClient_GetConsumerClientID_Call {
	_c.Call.Return(run)
	return _c
}

// GetMonitoringInfo provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringpClient) GetMonitoringInfo(ctx context.Context, in *types.QueryGetMonitoringInfoRequest, opts ...grpc.CallOption) (*types.QueryGetMonitoringInfoResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetMonitoringInfo")
	}

	var r0 *types.QueryGetMonitoringInfoResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetMonitoringInfoRequest, ...grpc.CallOption) (*types.QueryGetMonitoringInfoResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetMonitoringInfoRequest, ...grpc.CallOption) *types.QueryGetMonitoringInfoResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetMonitoringInfoResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetMonitoringInfoRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MonitoringpClient_GetMonitoringInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMonitoringInfo'
type MonitoringpClient_GetMonitoringInfo_Call struct {
	*mock.Call
}

// GetMonitoringInfo is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetMonitoringInfoRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringpClient_Expecter) GetMonitoringInfo(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringpClient_GetMonitoringInfo_Call {
	return &MonitoringpClient_GetMonitoringInfo_Call{Call: _e.mock.On("GetMonitoringInfo",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringpClient_GetMonitoringInfo_Call) Run(run func(ctx context.Context, in *types.QueryGetMonitoringInfoRequest, opts ...grpc.CallOption)) *MonitoringpClient_GetMonitoringInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetMonitoringInfoRequest), variadicArgs...)
	})
	return _c
}

func (_c *MonitoringpClient_GetMonitoringInfo_Call) Return(_a0 *types.QueryGetMonitoringInfoResponse, _a1 error) *MonitoringpClient_GetMonitoringInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringpClient_GetMonitoringInfo_Call) RunAndReturn(run func(context.Context, *types.QueryGetMonitoringInfoRequest, ...grpc.CallOption) (*types.QueryGetMonitoringInfoResponse, error)) *MonitoringpClient_GetMonitoringInfo_Call {
	_c.Call.Return(run)
	return _c
}

// Params provides a mock function with given fields: ctx, in, opts
func (_m *MonitoringpClient) Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error) {
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

// MonitoringpClient_Params_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Params'
type MonitoringpClient_Params_Call struct {
	*mock.Call
}

// Params is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryParamsRequest
//   - opts ...grpc.CallOption
func (_e *MonitoringpClient_Expecter) Params(ctx interface{}, in interface{}, opts ...interface{}) *MonitoringpClient_Params_Call {
	return &MonitoringpClient_Params_Call{Call: _e.mock.On("Params",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MonitoringpClient_Params_Call) Run(run func(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption)) *MonitoringpClient_Params_Call {
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

func (_c *MonitoringpClient_Params_Call) Return(_a0 *types.QueryParamsResponse, _a1 error) *MonitoringpClient_Params_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MonitoringpClient_Params_Call) RunAndReturn(run func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) (*types.QueryParamsResponse, error)) *MonitoringpClient_Params_Call {
	_c.Call.Return(run)
	return _c
}

// NewMonitoringpClient creates a new instance of MonitoringpClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMonitoringpClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MonitoringpClient {
	mock := &MonitoringpClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
