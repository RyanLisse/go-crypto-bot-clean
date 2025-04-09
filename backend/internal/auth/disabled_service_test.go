package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDisabledService_Authenticate(t *testing.T) {
	svc := &DisabledService{}
	req := httptest.NewRequest("GET", "/", nil)

	userData, err := svc.Authenticate(req)

	if userData != nil {
		t.Errorf("Expected nil user data, got %+v", userData)
	}

	if err == nil {
		t.Error("Expected error, got nil")
	}

	authErr, ok := err.(*AuthError)
	if !ok {
		t.Errorf("Expected AuthError, got %T", err)
	}

	if authErr.Type != ErrorTypeServiceUnavailable {
		t.Errorf("Expected error type %s, got %s", ErrorTypeServiceUnavailable, authErr.Type)
	}

	if authErr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, authErr.Code)
	}
}

func TestDisabledService_RequireRole(t *testing.T) {
	svc := &DisabledService{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware with the role requirement
	handler := svc.RequireRole("admin")(nextHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Check response body
	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response.Error.Type != ErrorTypeServiceUnavailable {
		t.Errorf("Expected error type %s, got %s", ErrorTypeServiceUnavailable, response.Error.Type)
	}

	// Check details
	details, ok := response.Error.Details.(map[string]interface{})
	if !ok {
		t.Errorf("Expected details to be a map, got %T", response.Error.Details)
	}

	if role, ok := details["requested_role"]; !ok || role != "admin" {
		t.Errorf("Expected requested_role to be 'admin', got %v", role)
	}
}

func TestDisabledService_RequirePermission(t *testing.T) {
	svc := &DisabledService{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware with the permission requirement
	handler := svc.RequirePermission("read:trades")(nextHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Check response body
	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response.Error.Type != ErrorTypeServiceUnavailable {
		t.Errorf("Expected error type %s, got %s", ErrorTypeServiceUnavailable, response.Error.Type)
	}

	// Check details
	details, ok := response.Error.Details.(map[string]interface{})
	if !ok {
		t.Errorf("Expected details to be a map, got %T", response.Error.Details)
	}

	if permission, ok := details["requested_permission"]; !ok || permission != "read:trades" {
		t.Errorf("Expected requested_permission to be 'read:trades', got %v", permission)
	}
}

func TestDisabledService_RequireAnyPermission(t *testing.T) {
	svc := &DisabledService{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware with the permission requirements
	handler := svc.RequireAnyPermission("read:trades", "create:trades")(nextHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Check response body
	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response.Error.Type != ErrorTypeServiceUnavailable {
		t.Errorf("Expected error type %s, got %s", ErrorTypeServiceUnavailable, response.Error.Type)
	}

	// Check details
	details, ok := response.Error.Details.(map[string]interface{})
	if !ok {
		t.Errorf("Expected details to be a map, got %T", response.Error.Details)
	}

	permissions, ok := details["requested_permissions"].([]interface{})
	if !ok {
		t.Errorf("Expected requested_permissions to be an array, got %T", details["requested_permissions"])
	}

	if len(permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions))
	}

	if permissions[0] != "read:trades" || permissions[1] != "create:trades" {
		t.Errorf("Expected permissions to be ['read:trades', 'create:trades'], got %v", permissions)
	}
}

func TestDisabledService_RequireAllPermissions(t *testing.T) {
	svc := &DisabledService{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware with the permission requirements
	handler := svc.RequireAllPermissions("read:trades", "create:trades")(nextHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Check response body
	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response.Error.Type != ErrorTypeServiceUnavailable {
		t.Errorf("Expected error type %s, got %s", ErrorTypeServiceUnavailable, response.Error.Type)
	}

	// Check details
	details, ok := response.Error.Details.(map[string]interface{})
	if !ok {
		t.Errorf("Expected details to be a map, got %T", response.Error.Details)
	}

	permissions, ok := details["required_permissions"].([]interface{})
	if !ok {
		t.Errorf("Expected required_permissions to be an array, got %T", details["required_permissions"])
	}

	if len(permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions))
	}

	if permissions[0] != "read:trades" || permissions[1] != "create:trades" {
		t.Errorf("Expected permissions to be ['read:trades', 'create:trades'], got %v", permissions)
	}
}

func TestDisabledService_AuthMiddleware(t *testing.T) {
	svc := &DisabledService{}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Create a test handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create middleware
	handler := svc.AuthMiddleware(nextHandler)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check status code - this should pass through
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	if rr.Body.String() != "success" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "success")
	}
}
