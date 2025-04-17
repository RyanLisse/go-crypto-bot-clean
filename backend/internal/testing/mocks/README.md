# Testing Mocks

This package contains isolated mocks for interfaces, designed to avoid import cycles and conflicts when testing components. These mocks implement interfaces similar to those in the domain but without directly importing from packages with conflicts.

## Why Isolated Mocks?

Sometimes projects have interface conflicts, such as:
- Duplicate interface declarations (e.g., `NewCoinUseCase` declared in multiple files)
- Inconsistent types (e.g., `Status` vs `CoinStatus`) 
- Refactoring in progress

This package provides a solution that allows you to continue testing components without directly importing from conflicting packages.

## Available Mocks

### MockNewCoinUseCase

A mock implementation of `NewCoinUseCase` that provides mocked methods for:
- `DetectNewCoins()`
- `UpdateCoinStatus()`
- `GetCoinDetails()`
- `ListNewCoins()`
- `GetRecentTradableCoins()`
- `SubscribeToEvents()`
- `UnsubscribeFromEvents()`

## Basic Usage

```go
func TestWithMock(t *testing.T) {
    // Create a mock
    mockNewCoinUC := &mocks.MockNewCoinUseCase{}

    // Setup expectations
    mockNewCoinUC.On("DetectNewCoins").Return(nil)

    // Call method
    err := mockNewCoinUC.DetectNewCoins()

    // Assert expectations
    assert.NoError(t, err)
    mockNewCoinUC.AssertExpectations(t)
}
```

## Using with Domain Interfaces

For components that expect the domain interfaces, use the adapter:

```go
func TestComponentWithDomainInterface(t *testing.T) {
    // Create a mock with adapter
    mockAdapter := mocks.NewMockNewCoinUseCase()

    // Create component that expects domain interface
    component := &YourComponent{
        coinUC: mockAdapter,  // Adapter implements domain interface
    }

    // Setup expectations on the underlying mock
    mockAdapter.Mock.On("DetectNewCoins").Return(nil)

    // Call component method that uses the domain interface
    result := component.DoSomething()

    // Assert expectations were met
    mockAdapter.Mock.AssertExpectations(t)
}
```

## Handling Events and Callbacks

For methods like `SubscribeToEvents` that use callbacks:

```go
func TestEventSubscription(t *testing.T) {
    // Create mock adapter
    mockAdapter := mocks.NewMockNewCoinUseCase()
    
    // Track if callback was called
    callbackCalled := false
    
    // Capture the registered callback
    var capturedCallback func(*mocks.NewCoinEvent)
    mockAdapter.Mock.On("SubscribeToEvents", mock.AnythingOfType("func(*mocks.NewCoinEvent)")).
        Run(func(args mock.Arguments) {
            capturedCallback = args.Get(0).(func(*mocks.NewCoinEvent))
        }).
        Return(nil)
    
    // Subscribe with callback
    err := mockAdapter.SubscribeToEvents(func(event *model.NewCoinEvent) {
        callbackCalled = true
    })
    
    // Simulate an event
    mockEvent := &mocks.NewCoinEvent{
        ID:        "event1",
        CoinID:    "BTC-USDT",
        EventType: "status_change",
    }
    capturedCallback(mockEvent)
    
    // Verify callback was called
    assert.True(t, callbackCalled)
}
```

## Examples

See the `examples/` directory for complete working examples of how to use these mocks in tests.

## Adding New Mocks

When adding new mocks to this package:

1. Create the mock in its own file (e.g., `mock_your_interface.go`)
2. Define any necessary types in the same file
3. Add unit tests for the mock
4. If needed, create an adapter to interface with domain types
5. Add examples to the examples directory 