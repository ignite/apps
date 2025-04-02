// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	"context"

	"github.com/ignite/network/x/profile/types"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// ProfileClient is an autogenerated mock type for the ProfileClient type
type ProfileClient struct {
	mock.Mock
}

type ProfileClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ProfileClient) EXPECT() *ProfileClient_Expecter {
	return &ProfileClient_Expecter{mock: &_m.Mock}
}

// GetCoordinator provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) GetCoordinator(ctx context.Context, in *types.QueryGetCoordinatorRequest, opts ...grpc.CallOption) (*types.QueryGetCoordinatorResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCoordinator")
	}

	var r0 *types.QueryGetCoordinatorResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetCoordinatorRequest, ...grpc.CallOption) (*types.QueryGetCoordinatorResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetCoordinatorRequest, ...grpc.CallOption) *types.QueryGetCoordinatorResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetCoordinatorResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetCoordinatorRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileClient_GetCoordinator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCoordinator'
type ProfileClient_GetCoordinator_Call struct {
	*mock.Call
}

// GetCoordinator is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetCoordinatorRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) GetCoordinator(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_GetCoordinator_Call {
	return &ProfileClient_GetCoordinator_Call{Call: _e.mock.On("GetCoordinator",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_GetCoordinator_Call) Run(run func(ctx context.Context, in *types.QueryGetCoordinatorRequest, opts ...grpc.CallOption)) *ProfileClient_GetCoordinator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetCoordinatorRequest), variadicArgs...)
	})
	return _c
}

func (_c *ProfileClient_GetCoordinator_Call) Return(_a0 *types.QueryGetCoordinatorResponse, _a1 error) *ProfileClient_GetCoordinator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_GetCoordinator_Call) RunAndReturn(run func(context.Context, *types.QueryGetCoordinatorRequest, ...grpc.CallOption) (*types.QueryGetCoordinatorResponse, error)) *ProfileClient_GetCoordinator_Call {
	_c.Call.Return(run)
	return _c
}

// GetCoordinatorByAddress provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) GetCoordinatorByAddress(ctx context.Context, in *types.QueryGetCoordinatorByAddressRequest, opts ...grpc.CallOption) (*types.QueryGetCoordinatorByAddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetCoordinatorByAddress")
	}

	var r0 *types.QueryGetCoordinatorByAddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetCoordinatorByAddressRequest, ...grpc.CallOption) (*types.QueryGetCoordinatorByAddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetCoordinatorByAddressRequest, ...grpc.CallOption) *types.QueryGetCoordinatorByAddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetCoordinatorByAddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetCoordinatorByAddressRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileClient_GetCoordinatorByAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCoordinatorByAddress'
type ProfileClient_GetCoordinatorByAddress_Call struct {
	*mock.Call
}

// GetCoordinatorByAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetCoordinatorByAddressRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) GetCoordinatorByAddress(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_GetCoordinatorByAddress_Call {
	return &ProfileClient_GetCoordinatorByAddress_Call{Call: _e.mock.On("GetCoordinatorByAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_GetCoordinatorByAddress_Call) Run(run func(ctx context.Context, in *types.QueryGetCoordinatorByAddressRequest, opts ...grpc.CallOption)) *ProfileClient_GetCoordinatorByAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetCoordinatorByAddressRequest), variadicArgs...)
	})
	return _c
}

func (_c *ProfileClient_GetCoordinatorByAddress_Call) Return(_a0 *types.QueryGetCoordinatorByAddressResponse, _a1 error) *ProfileClient_GetCoordinatorByAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_GetCoordinatorByAddress_Call) RunAndReturn(run func(context.Context, *types.QueryGetCoordinatorByAddressRequest, ...grpc.CallOption) (*types.QueryGetCoordinatorByAddressResponse, error)) *ProfileClient_GetCoordinatorByAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetValidator provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) GetValidator(ctx context.Context, in *types.QueryGetValidatorRequest, opts ...grpc.CallOption) (*types.QueryGetValidatorResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetValidator")
	}

	var r0 *types.QueryGetValidatorResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetValidatorRequest, ...grpc.CallOption) (*types.QueryGetValidatorResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetValidatorRequest, ...grpc.CallOption) *types.QueryGetValidatorResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetValidatorResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetValidatorRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileClient_GetValidator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValidator'
type ProfileClient_GetValidator_Call struct {
	*mock.Call
}

// GetValidator is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetValidatorRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) GetValidator(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_GetValidator_Call {
	return &ProfileClient_GetValidator_Call{Call: _e.mock.On("GetValidator",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_GetValidator_Call) Run(run func(ctx context.Context, in *types.QueryGetValidatorRequest, opts ...grpc.CallOption)) *ProfileClient_GetValidator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetValidatorRequest), variadicArgs...)
	})
	return _c
}

func (_c *ProfileClient_GetValidator_Call) Return(_a0 *types.QueryGetValidatorResponse, _a1 error) *ProfileClient_GetValidator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_GetValidator_Call) RunAndReturn(run func(context.Context, *types.QueryGetValidatorRequest, ...grpc.CallOption) (*types.QueryGetValidatorResponse, error)) *ProfileClient_GetValidator_Call {
	_c.Call.Return(run)
	return _c
}

// GetValidatorByOperatorAddress provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) GetValidatorByOperatorAddress(ctx context.Context, in *types.QueryGetValidatorByOperatorAddressRequest, opts ...grpc.CallOption) (*types.QueryGetValidatorByOperatorAddressResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetValidatorByOperatorAddress")
	}

	var r0 *types.QueryGetValidatorByOperatorAddressResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetValidatorByOperatorAddressRequest, ...grpc.CallOption) (*types.QueryGetValidatorByOperatorAddressResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetValidatorByOperatorAddressRequest, ...grpc.CallOption) *types.QueryGetValidatorByOperatorAddressResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetValidatorByOperatorAddressResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetValidatorByOperatorAddressRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileClient_GetValidatorByOperatorAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValidatorByOperatorAddress'
type ProfileClient_GetValidatorByOperatorAddress_Call struct {
	*mock.Call
}

// GetValidatorByOperatorAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryGetValidatorByOperatorAddressRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) GetValidatorByOperatorAddress(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_GetValidatorByOperatorAddress_Call {
	return &ProfileClient_GetValidatorByOperatorAddress_Call{Call: _e.mock.On("GetValidatorByOperatorAddress",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_GetValidatorByOperatorAddress_Call) Run(run func(ctx context.Context, in *types.QueryGetValidatorByOperatorAddressRequest, opts ...grpc.CallOption)) *ProfileClient_GetValidatorByOperatorAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryGetValidatorByOperatorAddressRequest), variadicArgs...)
	})
	return _c
}

func (_c *ProfileClient_GetValidatorByOperatorAddress_Call) Return(_a0 *types.QueryGetValidatorByOperatorAddressResponse, _a1 error) *ProfileClient_GetValidatorByOperatorAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_GetValidatorByOperatorAddress_Call) RunAndReturn(run func(context.Context, *types.QueryGetValidatorByOperatorAddressRequest, ...grpc.CallOption) (*types.QueryGetValidatorByOperatorAddressResponse, error)) *ProfileClient_GetValidatorByOperatorAddress_Call {
	_c.Call.Return(run)
	return _c
}

// ListCoordinator provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) ListCoordinator(ctx context.Context, in *types.QueryAllCoordinatorRequest, opts ...grpc.CallOption) (*types.QueryAllCoordinatorResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListCoordinator")
	}

	var r0 *types.QueryAllCoordinatorResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllCoordinatorRequest, ...grpc.CallOption) (*types.QueryAllCoordinatorResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllCoordinatorRequest, ...grpc.CallOption) *types.QueryAllCoordinatorResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllCoordinatorResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllCoordinatorRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileClient_ListCoordinator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListCoordinator'
type ProfileClient_ListCoordinator_Call struct {
	*mock.Call
}

// ListCoordinator is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryAllCoordinatorRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) ListCoordinator(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_ListCoordinator_Call {
	return &ProfileClient_ListCoordinator_Call{Call: _e.mock.On("ListCoordinator",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_ListCoordinator_Call) Run(run func(ctx context.Context, in *types.QueryAllCoordinatorRequest, opts ...grpc.CallOption)) *ProfileClient_ListCoordinator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryAllCoordinatorRequest), variadicArgs...)
	})
	return _c
}

func (_c *ProfileClient_ListCoordinator_Call) Return(_a0 *types.QueryAllCoordinatorResponse, _a1 error) *ProfileClient_ListCoordinator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_ListCoordinator_Call) RunAndReturn(run func(context.Context, *types.QueryAllCoordinatorRequest, ...grpc.CallOption) (*types.QueryAllCoordinatorResponse, error)) *ProfileClient_ListCoordinator_Call {
	_c.Call.Return(run)
	return _c
}

// ListValidator provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) ListValidator(ctx context.Context, in *types.QueryAllValidatorRequest, opts ...grpc.CallOption) (*types.QueryAllValidatorResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListValidator")
	}

	var r0 *types.QueryAllValidatorResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllValidatorRequest, ...grpc.CallOption) (*types.QueryAllValidatorResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllValidatorRequest, ...grpc.CallOption) *types.QueryAllValidatorResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllValidatorResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllValidatorRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileClient_ListValidator_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListValidator'
type ProfileClient_ListValidator_Call struct {
	*mock.Call
}

// ListValidator is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryAllValidatorRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) ListValidator(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_ListValidator_Call {
	return &ProfileClient_ListValidator_Call{Call: _e.mock.On("ListValidator",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_ListValidator_Call) Run(run func(ctx context.Context, in *types.QueryAllValidatorRequest, opts ...grpc.CallOption)) *ProfileClient_ListValidator_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*types.QueryAllValidatorRequest), variadicArgs...)
	})
	return _c
}

func (_c *ProfileClient_ListValidator_Call) Return(_a0 *types.QueryAllValidatorResponse, _a1 error) *ProfileClient_ListValidator_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_ListValidator_Call) RunAndReturn(run func(context.Context, *types.QueryAllValidatorRequest, ...grpc.CallOption) (*types.QueryAllValidatorResponse, error)) *ProfileClient_ListValidator_Call {
	_c.Call.Return(run)
	return _c
}

// Params provides a mock function with given fields: ctx, in, opts
func (_m *ProfileClient) Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error) {
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

// ProfileClient_Params_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Params'
type ProfileClient_Params_Call struct {
	*mock.Call
}

// Params is a helper method to define mock.On call
//   - ctx context.Context
//   - in *types.QueryParamsRequest
//   - opts ...grpc.CallOption
func (_e *ProfileClient_Expecter) Params(ctx interface{}, in interface{}, opts ...interface{}) *ProfileClient_Params_Call {
	return &ProfileClient_Params_Call{Call: _e.mock.On("Params",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *ProfileClient_Params_Call) Run(run func(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption)) *ProfileClient_Params_Call {
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

func (_c *ProfileClient_Params_Call) Return(_a0 *types.QueryParamsResponse, _a1 error) *ProfileClient_Params_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileClient_Params_Call) RunAndReturn(run func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) (*types.QueryParamsResponse, error)) *ProfileClient_Params_Call {
	_c.Call.Return(run)
	return _c
}

// NewProfileClient creates a new instance of ProfileClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProfileClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProfileClient {
	mock := &ProfileClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
