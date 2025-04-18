package mocks

package mocks

import (
	context "context"

	model "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// RiskMetricsRepository is an autogenerated mock type for the RiskMetricsRepository type
type RiskMetricsRepository struct {
	mock.Mock
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *RiskMetricsRepository) GetByUserID(ctx context.Context, userID string) (*model.RiskMetrics, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 *model.RiskMetrics
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.RiskMetrics, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.RiskMetrics); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RiskMetrics)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistorical provides a mock function with given fields: ctx, userID, from, to, interval
func (_m *RiskMetricsRepository) GetHistorical(ctx context.Context, userID string, from time.Time, to time.Time, interval string) ([]*model.RiskMetrics, error) {
	ret := _m.Called(ctx, userID, from, to, interval)

	if len(ret) == 0 {
		panic("no return value specified for GetHistorical")
	}

	var r0 []*model.RiskMetrics
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time, string) ([]*model.RiskMetrics, error)); ok {
		return rf(ctx, userID, from, to, interval)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time, time.Time, string) []*model.RiskMetrics); ok {
		r0 = rf(ctx, userID, from, to, interval)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.RiskMetrics)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time, time.Time, string) error); ok {
		r1 = rf(ctx, userID, from, to, interval)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, metrics
func (_m *RiskMetricsRepository) Save(ctx context.Context, metrics *model.RiskMetrics) error {
	ret := _m.Called(ctx, metrics)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.RiskMetrics) error); ok {
		r0 = rf(ctx, metrics)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRiskMetricsRepository creates a new instance of RiskMetricsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRiskMetricsRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *RiskMetricsRepository {
	mock := &RiskMetricsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
