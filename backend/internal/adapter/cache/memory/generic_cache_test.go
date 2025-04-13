package memory

import (
	"errors"
	"testing"
	"time"
)

type testStruct struct {
	ID   string
	Name string
}

func TestGenericCache_Get(t *testing.T) {
	cache := NewGenericCache[testStruct](time.Second)

	// Initially the cache should be empty
	value, found := cache.Get()
	if found {
		t.Errorf("Expected empty cache to return found=false, got found=true")
	}
	if value != nil {
		t.Errorf("Expected empty cache to return nil value, got %v", value)
	}

	// Set a value and verify we can get it
	testValue := &testStruct{ID: "1", Name: "Test"}
	cache.Set(testValue)

	value, found = cache.Get()
	if !found {
		t.Errorf("Expected value to be found in cache, got found=false")
	}
	if value == nil {
		t.Errorf("Expected non-nil value from cache")
	} else if value.ID != testValue.ID || value.Name != testValue.Name {
		t.Errorf("Expected value %v, got %v", testValue, value)
	}

	// Test expiration
	cache = NewGenericCache[testStruct](10 * time.Millisecond)
	cache.Set(testValue)

	time.Sleep(20 * time.Millisecond)

	value, found = cache.Get()
	if found {
		t.Errorf("Expected value to be expired, got found=true")
	}
	if value != nil {
		t.Errorf("Expected nil value for expired entry, got %v", value)
	}
}

func TestGenericCache_GetOrSet(t *testing.T) {
	cache := NewGenericCache[testStruct](time.Second)

	// Test with successful fetch function
	fetchCalled := false
	fetchFn := func() (*testStruct, error) {
		fetchCalled = true
		return &testStruct{ID: "1", Name: "Test"}, nil
	}

	value, err := cache.GetOrSet(fetchFn)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value == nil {
		t.Errorf("Expected non-nil value")
	} else if value.ID != "1" || value.Name != "Test" {
		t.Errorf("Expected {ID:1, Name:Test}, got %v", value)
	}
	if !fetchCalled {
		t.Errorf("Expected fetch function to be called")
	}

	// Second call should use cached value
	fetchCalled = false
	value, err = cache.GetOrSet(fetchFn)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value == nil {
		t.Errorf("Expected non-nil value")
	}
	if fetchCalled {
		t.Errorf("Expected fetch function not to be called on cache hit")
	}

	// Test with failing fetch function
	cache = NewGenericCache[testStruct](time.Second)
	failingFetchFn := func() (*testStruct, error) {
		return nil, errors.New("fetch error")
	}

	value, err = cache.GetOrSet(failingFetchFn)
	if err == nil {
		t.Errorf("Expected error to be returned from failing fetch function")
	}
	if value != nil {
		t.Errorf("Expected nil value on error, got %v", value)
	}
}

func TestGenericCache_Invalidate(t *testing.T) {
	cache := NewGenericCache[testStruct](time.Minute)

	// Set a value
	testValue := &testStruct{ID: "1", Name: "Test"}
	cache.Set(testValue)

	// Verify it's there
	value, found := cache.Get()
	if !found || value == nil {
		t.Errorf("Expected value to be in cache before invalidation")
	}

	// Invalidate and verify it's gone
	cache.Invalidate()

	value, found = cache.Get()
	if found {
		t.Errorf("Expected value not to be found after invalidation")
	}
	if value != nil {
		t.Errorf("Expected nil value after invalidation, got %v", value)
	}
}

func TestGenericCache_UpdateTTL(t *testing.T) {
	cache := NewGenericCache[testStruct](100 * time.Millisecond)

	// Set a value
	testValue := &testStruct{ID: "1", Name: "Test"}
	cache.Set(testValue)

	// Update TTL to longer duration
	cache.UpdateTTL(10 * time.Minute)

	// Sleep past original TTL
	time.Sleep(200 * time.Millisecond)

	// Value should still be present due to extended TTL
	value, found := cache.Get()
	if !found {
		t.Errorf("Expected value to still be in cache after TTL extension")
	}
	if value == nil {
		t.Errorf("Expected non-nil value after TTL extension")
	}

	// Test with invalid TTL (should not change existing TTL)
	cache = NewGenericCache[testStruct](time.Minute)
	cache.Set(testValue)
	cache.UpdateTTL(-10 * time.Second)

	// Verify original TTL is still in effect
	// This is harder to test directly, but we can check that the value is still there
	value, found = cache.Get()
	if !found || value == nil {
		t.Errorf("Expected value to still be in cache after invalid TTL update")
	}
}
