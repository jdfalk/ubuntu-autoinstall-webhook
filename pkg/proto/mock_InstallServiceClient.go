// Code generated by mockery v2.53.0. DO NOT EDIT.

package proto

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// MockInstallServiceClient is an autogenerated mock type for the InstallServiceClient type
type MockInstallServiceClient struct {
	mock.Mock
}

type MockInstallServiceClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInstallServiceClient) EXPECT() *MockInstallServiceClient_Expecter {
	return &MockInstallServiceClient_Expecter{mock: &_m.Mock}
}

// ReportStatus provides a mock function with given fields: ctx, in, opts
func (_m *MockInstallServiceClient) ReportStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ReportStatus")
	}

	var r0 *StatusResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *StatusRequest, ...grpc.CallOption) (*StatusResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *StatusRequest, ...grpc.CallOption) *StatusResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*StatusResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *StatusRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInstallServiceClient_ReportStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReportStatus'
type MockInstallServiceClient_ReportStatus_Call struct {
	*mock.Call
}

// ReportStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - in *StatusRequest
//   - opts ...grpc.CallOption
func (_e *MockInstallServiceClient_Expecter) ReportStatus(ctx interface{}, in interface{}, opts ...interface{}) *MockInstallServiceClient_ReportStatus_Call {
	return &MockInstallServiceClient_ReportStatus_Call{Call: _e.mock.On("ReportStatus",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockInstallServiceClient_ReportStatus_Call) Run(run func(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption)) *MockInstallServiceClient_ReportStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*StatusRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockInstallServiceClient_ReportStatus_Call) Return(_a0 *StatusResponse, _a1 error) *MockInstallServiceClient_ReportStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInstallServiceClient_ReportStatus_Call) RunAndReturn(run func(context.Context, *StatusRequest, ...grpc.CallOption) (*StatusResponse, error)) *MockInstallServiceClient_ReportStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockInstallServiceClient creates a new instance of MockInstallServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInstallServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInstallServiceClient {
	mock := &MockInstallServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
