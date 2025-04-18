package mocks

package mocks

import (
	event "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/event"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// DomainEvent is an autogenerated mock type for the DomainEvent type
type DomainEvent struct {
	mock.Mock
}

// AggregateID provides a mock function with no fields
func (_m *DomainEvent) AggregateID() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AggregateID")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OccurredAt provides a mock function with no fields
func (_m *DomainEvent) OccurredAt() time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for OccurredAt")
	}

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// Type provides a mock function with no fields
func (_m *DomainEvent) Type() event.EventType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 event.EventType
	if rf, ok := ret.Get(0).(func() event.EventType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(event.EventType)
	}

	return r0
}

// NewDomainEvent creates a new instance of DomainEvent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDomainEvent(t interface {
	mock.TestingT
	Cleanup(func())
}) *DomainEvent {
	mock := &DomainEvent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
