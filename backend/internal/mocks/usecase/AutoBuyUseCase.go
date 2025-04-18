package mocks

package mocks

import (
	context "context"

	model "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// AutoBuyUseCase is an autogenerated mock type for the AutoBuyUseCase type
type AutoBuyUseCase struct {
	mock.Mock
}

// CreateAutoRule provides a mock function with given fields: ctx, userID, rule
func (_m *AutoBuyUseCase) CreateAutoRule(ctx context.Context, userID string, rule *model.AutoBuyRule) error {
	ret := _m.Called(ctx, userID, rule)

	if len(ret) == 0 {
		panic("no return value specified for CreateAutoRule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.AutoBuyRule) error); ok {
		r0 = rf(ctx, userID, rule)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAutoRule provides a mock function with given fields: ctx, ruleID
func (_m *AutoBuyUseCase) DeleteAutoRule(ctx context.Context, ruleID string) error {
	ret := _m.Called(ctx, ruleID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAutoRule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, ruleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EvaluateRule provides a mock function with given fields: ctx, ruleID
func (_m *AutoBuyUseCase) EvaluateRule(ctx context.Context, ruleID string) (*model.Order, error) {
	ret := _m.Called(ctx, ruleID)

	if len(ret) == 0 {
		panic("no return value specified for EvaluateRule")
	}

	var r0 *model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Order, error)); ok {
		return rf(ctx, ruleID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Order); ok {
		r0 = rf(ctx, ruleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ruleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EvaluateRules provides a mock function with given fields: ctx
func (_m *AutoBuyUseCase) EvaluateRules(ctx context.Context) ([]*model.Order, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for EvaluateRules")
	}

	var r0 []*model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.Order, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Order); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAutoRuleByID provides a mock function with given fields: ctx, ruleID
func (_m *AutoBuyUseCase) GetAutoRuleByID(ctx context.Context, ruleID string) (*model.AutoBuyRule, error) {
	ret := _m.Called(ctx, ruleID)

	if len(ret) == 0 {
		panic("no return value specified for GetAutoRuleByID")
	}

	var r0 *model.AutoBuyRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.AutoBuyRule, error)); ok {
		return rf(ctx, ruleID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.AutoBuyRule); ok {
		r0 = rf(ctx, ruleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AutoBuyRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ruleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAutoRulesBySymbol provides a mock function with given fields: ctx, symbol
func (_m *AutoBuyUseCase) GetAutoRulesBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error) {
	ret := _m.Called(ctx, symbol)

	if len(ret) == 0 {
		panic("no return value specified for GetAutoRulesBySymbol")
	}

	var r0 []*model.AutoBuyRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.AutoBuyRule, error)); ok {
		return rf(ctx, symbol)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.AutoBuyRule); ok {
		r0 = rf(ctx, symbol)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AutoBuyRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, symbol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAutoRulesByUser provides a mock function with given fields: ctx, userID
func (_m *AutoBuyUseCase) GetAutoRulesByUser(ctx context.Context, userID string) ([]*model.AutoBuyRule, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetAutoRulesByUser")
	}

	var r0 []*model.AutoBuyRule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.AutoBuyRule, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.AutoBuyRule); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AutoBuyRule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionHistory provides a mock function with given fields: ctx, userID, limit, offset
func (_m *AutoBuyUseCase) GetExecutionHistory(ctx context.Context, userID string, limit int, offset int) ([]*model.AutoBuyExecution, error) {
	ret := _m.Called(ctx, userID, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetExecutionHistory")
	}

	var r0 []*model.AutoBuyExecution
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]*model.AutoBuyExecution, error)); ok {
		return rf(ctx, userID, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*model.AutoBuyExecution); ok {
		r0 = rf(ctx, userID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AutoBuyExecution)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, userID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAutoRule provides a mock function with given fields: ctx, rule
func (_m *AutoBuyUseCase) UpdateAutoRule(ctx context.Context, rule *model.AutoBuyRule) error {
	ret := _m.Called(ctx, rule)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAutoRule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AutoBuyRule) error); ok {
		r0 = rf(ctx, rule)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAutoBuyUseCase creates a new instance of AutoBuyUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAutoBuyUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *AutoBuyUseCase {
	mock := &AutoBuyUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
