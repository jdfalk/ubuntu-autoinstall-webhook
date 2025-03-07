// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	configuration "github.com/jdfalk/ubuntu-autoinstall-webhook/internal/configuration"
	mock "github.com/stretchr/testify/mock"
)

// ConfigService is an autogenerated mock type for the ConfigService type
type ConfigService struct {
	mock.Mock
}

type ConfigService_Expecter struct {
	mock *mock.Mock
}

func (_m *ConfigService) EXPECT() *ConfigService_Expecter {
	return &ConfigService_Expecter{mock: &_m.Mock}
}

// GenerateTemplates provides a mock function with given fields: data
func (_m *ConfigService) GenerateTemplates(data interface{}) (map[string][]byte, error) {
	ret := _m.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for GenerateTemplates")
	}

	var r0 map[string][]byte
	var r1 error
	if rf, ok := ret.Get(0).(func(interface{}) (map[string][]byte, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func(interface{}) map[string][]byte); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConfigService_GenerateTemplates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateTemplates'
type ConfigService_GenerateTemplates_Call struct {
	*mock.Call
}

// GenerateTemplates is a helper method to define mock.On call
//   - data interface{}
func (_e *ConfigService_Expecter) GenerateTemplates(data interface{}) *ConfigService_GenerateTemplates_Call {
	return &ConfigService_GenerateTemplates_Call{Call: _e.mock.On("GenerateTemplates", data)}
}

func (_c *ConfigService_GenerateTemplates_Call) Run(run func(data interface{})) *ConfigService_GenerateTemplates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *ConfigService_GenerateTemplates_Call) Return(_a0 map[string][]byte, _a1 error) *ConfigService_GenerateTemplates_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ConfigService_GenerateTemplates_Call) RunAndReturn(run func(interface{}) (map[string][]byte, error)) *ConfigService_GenerateTemplates_Call {
	_c.Call.Return(run)
	return _c
}

// LoadConfig provides a mock function with no fields
func (_m *ConfigService) LoadConfig() (configuration.Config, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for LoadConfig")
	}

	var r0 configuration.Config
	var r1 error
	if rf, ok := ret.Get(0).(func() (configuration.Config, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() configuration.Config); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(configuration.Config)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConfigService_LoadConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadConfig'
type ConfigService_LoadConfig_Call struct {
	*mock.Call
}

// LoadConfig is a helper method to define mock.On call
func (_e *ConfigService_Expecter) LoadConfig() *ConfigService_LoadConfig_Call {
	return &ConfigService_LoadConfig_Call{Call: _e.mock.On("LoadConfig")}
}

func (_c *ConfigService_LoadConfig_Call) Run(run func()) *ConfigService_LoadConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ConfigService_LoadConfig_Call) Return(_a0 configuration.Config, _a1 error) *ConfigService_LoadConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ConfigService_LoadConfig_Call) RunAndReturn(run func() (configuration.Config, error)) *ConfigService_LoadConfig_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateConfig provides a mock function with given fields: cfg
func (_m *ConfigService) ValidateConfig(cfg configuration.Config) error {
	ret := _m.Called(cfg)

	if len(ret) == 0 {
		panic("no return value specified for ValidateConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(configuration.Config) error); ok {
		r0 = rf(cfg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConfigService_ValidateConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateConfig'
type ConfigService_ValidateConfig_Call struct {
	*mock.Call
}

// ValidateConfig is a helper method to define mock.On call
//   - cfg configuration.Config
func (_e *ConfigService_Expecter) ValidateConfig(cfg interface{}) *ConfigService_ValidateConfig_Call {
	return &ConfigService_ValidateConfig_Call{Call: _e.mock.On("ValidateConfig", cfg)}
}

func (_c *ConfigService_ValidateConfig_Call) Run(run func(cfg configuration.Config)) *ConfigService_ValidateConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(configuration.Config))
	})
	return _c
}

func (_c *ConfigService_ValidateConfig_Call) Return(_a0 error) *ConfigService_ValidateConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ConfigService_ValidateConfig_Call) RunAndReturn(run func(configuration.Config) error) *ConfigService_ValidateConfig_Call {
	_c.Call.Return(run)
	return _c
}

// WatchConfigUpdates provides a mock function with no fields
func (_m *ConfigService) WatchConfigUpdates() (<-chan configuration.Config, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for WatchConfigUpdates")
	}

	var r0 <-chan configuration.Config
	var r1 error
	if rf, ok := ret.Get(0).(func() (<-chan configuration.Config, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() <-chan configuration.Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan configuration.Config)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConfigService_WatchConfigUpdates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WatchConfigUpdates'
type ConfigService_WatchConfigUpdates_Call struct {
	*mock.Call
}

// WatchConfigUpdates is a helper method to define mock.On call
func (_e *ConfigService_Expecter) WatchConfigUpdates() *ConfigService_WatchConfigUpdates_Call {
	return &ConfigService_WatchConfigUpdates_Call{Call: _e.mock.On("WatchConfigUpdates")}
}

func (_c *ConfigService_WatchConfigUpdates_Call) Run(run func()) *ConfigService_WatchConfigUpdates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ConfigService_WatchConfigUpdates_Call) Return(_a0 <-chan configuration.Config, _a1 error) *ConfigService_WatchConfigUpdates_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ConfigService_WatchConfigUpdates_Call) RunAndReturn(run func() (<-chan configuration.Config, error)) *ConfigService_WatchConfigUpdates_Call {
	_c.Call.Return(run)
	return _c
}

// NewConfigService creates a new instance of ConfigService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConfigService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConfigService {
	mock := &ConfigService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
