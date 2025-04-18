package mocks

package mocks

import (
	context "context"

	model "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// WalletRepository is an autogenerated mock type for the WalletRepository type
type WalletRepository struct {
	mock.Mock
}

// GetBalanceHistory provides a mock function with given fields: ctx, userID, asset, from, to
func (_m *WalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from time.Time, to time.Time) ([]*model.BalanceHistory, error) {
	ret := _m.Called(ctx, userID, asset, from, to)

	if len(ret) == 0 {
		panic("no return value specified for GetBalanceHistory")
	}

	var r0 []*model.BalanceHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.Asset, time.Time, time.Time) ([]*model.BalanceHistory, error)); ok {
		return rf(ctx, userID, asset, from, to)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, model.Asset, time.Time, time.Time) []*model.BalanceHistory); ok {
		r0 = rf(ctx, userID, asset, from, to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.BalanceHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, model.Asset, time.Time, time.Time) error); ok {
		r1 = rf(ctx, userID, asset, from, to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUserID provides a mock function with given fields: ctx, userID
func (_m *WalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 *model.Wallet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Wallet, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Wallet); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Wallet)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, wallet
func (_m *WalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	ret := _m.Called(ctx, wallet)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Wallet) error); ok {
		r0 = rf(ctx, wallet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveBalanceHistory provides a mock function with given fields: ctx, history
func (_m *WalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	ret := _m.Called(ctx, history)

	if len(ret) == 0 {
		panic("no return value specified for SaveBalanceHistory")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.BalanceHistory) error); ok {
		r0 = rf(ctx, history)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewWalletRepository creates a new instance of WalletRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWalletRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *WalletRepository {
	mock := &WalletRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
