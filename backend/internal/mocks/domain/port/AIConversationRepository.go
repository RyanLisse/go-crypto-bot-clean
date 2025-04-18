package mocks

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AIConversationRepository is an autogenerated mock type for the AIConversationRepository type
type AIConversationRepository struct {
	mock.Mock
}

// ClearConversation provides a mock function with given fields: ctx, userID
func (_m *AIConversationRepository) ClearConversation(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for ClearConversation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetConversation provides a mock function with given fields: ctx, userID, limit
func (_m *AIConversationRepository) GetConversation(ctx context.Context, userID string, limit int) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, userID, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetConversation")
	}

	var r0 []map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) ([]map[string]interface{}, error)); ok {
		return rf(ctx, userID, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) []map[string]interface{}); ok {
		r0 = rf(ctx, userID, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, userID, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveMessage provides a mock function with given fields: ctx, userID, message
func (_m *AIConversationRepository) SaveMessage(ctx context.Context, userID string, message map[string]interface{}) error {
	ret := _m.Called(ctx, userID, message)

	if len(ret) == 0 {
		panic("no return value specified for SaveMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) error); ok {
		r0 = rf(ctx, userID, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewAIConversationRepository creates a new instance of AIConversationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAIConversationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AIConversationRepository {
	mock := &AIConversationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
