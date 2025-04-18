package mocks

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StrategyRepository is an autogenerated mock type for the StrategyRepository type
type StrategyRepository struct {
	mock.Mock
}

// DeleteStrategy provides a mock function with given fields: ctx, strategyID
func (_m *StrategyRepository) DeleteStrategy(ctx context.Context, strategyID string) error {
	ret := _m.Called(ctx, strategyID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteStrategy")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, strategyID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetConfig provides a mock function with given fields: ctx, strategyID
func (_m *StrategyRepository) GetConfig(ctx context.Context, strategyID string) (map[string]interface{}, error) {
	ret := _m.Called(ctx, strategyID)

	if len(ret) == 0 {
		panic("no return value specified for GetConfig")
	}

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (map[string]interface{}, error)); ok {
		return rf(ctx, strategyID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]interface{}); ok {
		r0 = rf(ctx, strategyID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, strategyID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListStrategies provides a mock function with given fields: ctx
func (_m *StrategyRepository) ListStrategies(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListStrategies")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveConfig provides a mock function with given fields: ctx, strategyID, config
func (_m *StrategyRepository) SaveConfig(ctx context.Context, strategyID string, config map[string]interface{}) error {
	ret := _m.Called(ctx, strategyID, config)

	if len(ret) == 0 {
		panic("no return value specified for SaveConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) error); ok {
		r0 = rf(ctx, strategyID, config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStrategyRepository creates a new instance of StrategyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStrategyRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *StrategyRepository {
	mock := &StrategyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
