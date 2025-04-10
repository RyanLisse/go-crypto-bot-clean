package adapter

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewChiRouter(t *testing.T) {
	router := NewChiRouter()
	assert.NotNil(t, router)
	assert.NotNil(t, router.Router())
}

func TestChiRouterBasicRoutes(t *testing.T) {
	router := NewChiRouter()

	// Test GET route
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "get"})
	})

	// Test POST route
	router.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "post"})
	})

	// Test PUT route
	router.Put("/test", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "put"})
	})

	// Test DELETE route
	router.Delete("/test", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "delete"})
	})

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{"GET request", "GET", http.StatusOK, `{"success":true,"data":{"message":"get"}}`},
		{"POST request", "POST", http.StatusOK, `{"success":true,"data":{"message":"post"}}`},
		{"PUT request", "PUT", http.StatusOK, `{"success":true,"data":{"message":"put"}}`},
		{"DELETE request", "DELETE", http.StatusOK, `{"success":true,"data":{"message":"delete"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()
			router.Router().ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestChiRouter(t *testing.T) {
	// Test cases
	tests := []struct {
		name         string
		method       string
		path         string
		handler      http.HandlerFunc
		expectedCode int
		expectedBody string
	}{
		{
			name:   "Simple GET handler",
			method: "GET",
			path:   "/test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				RespondWithJSON(w, http.StatusOK, map[string]string{"message": "hello"})
			},
			expectedCode: http.StatusOK,
			expectedBody: `{"success":true,"data":{"message":"hello"}}`,
		},
		{
			name:   "Simple POST handler",
			method: "POST",
			path:   "/test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "created"})
			},
			expectedCode: http.StatusCreated,
			expectedBody: `{"success":true,"data":{"status":"created"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create router
			router := NewChiRouter()

			// Register handler
			switch tt.method {
			case "GET":
				router.Get(tt.path, tt.handler)
			case "POST":
				router.Post(tt.path, tt.handler)
			}

			// Create request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Serve request
			router.Router().ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Check response body
			if tt.expectedBody != "" {
				var expected, actual interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expected)
				assert.NoError(t, err)
				err = json.Unmarshal(w.Body.Bytes(), &actual)
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestChiRouterWithParams(t *testing.T) {
	router := NewChiRouter()

	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := GetParam(r, "id")
		RespondWithJSON(w, http.StatusOK, map[string]string{"id": id})
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"success":true,"data":{"id":"123"}}`, w.Body.String())
}

func TestChiRouterErrorHandling(t *testing.T) {
	router := NewChiRouter()

	router.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		RespondWithError(w, http.StatusBadRequest, "test error")
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"success":false,"error":"test error"}`, w.Body.String())
}

func TestChiRouterMiddleware(t *testing.T) {
	router := NewChiRouter()

	// Add custom middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "middleware")
			next.ServeHTTP(w, r)
		})
	})

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "middleware", w.Header().Get("X-Test"))
	assert.JSONEq(t, `{"success":true,"data":{"message":"test"}}`, w.Body.String())
}

func TestChiRouterWithContext(t *testing.T) {
	router := NewChiRouter()

	type contextKey string
	testKey := contextKey("test")

	router.WithContext(testKey, "test-value")
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		value := r.Context().Value(testKey).(string)
		RespondWithJSON(w, http.StatusOK, map[string]string{"value": value})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"success":true,"data":{"value":"test-value"}}`, w.Body.String())
}

func TestChiRouterAdditionalMethods(t *testing.T) {
	router := NewChiRouter()

	// Test Method function
	router.Method("PATCH", "/test", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"method": "patch"})
	})

	// Test Handle function
	router.Handle("/handle", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]string{"type": "handler"})
	}))

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			"PATCH request",
			"PATCH",
			"/test",
			http.StatusOK,
			`{"success":true,"data":{"method":"patch"}}`,
		},
		{
			"Handle function",
			"GET",
			"/handle",
			http.StatusOK,
			`{"success":true,"data":{"type":"handler"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			router.Router().ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestGetBody(t *testing.T) {
	router := NewChiRouter()

	router.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		body, err := GetBody(r)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		var data map[string]string
		if err := json.Unmarshal(body, &data); err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, data)
	})

	payload := `{"message":"test body"}`
	req := httptest.NewRequest("POST", "/test", strings.NewReader(payload))
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"success":true,"data":{"message":"test body"}}`, w.Body.String())
}

func TestGetQuery(t *testing.T) {
	router := NewChiRouter()

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		value := GetQuery(r, "key")
		RespondWithJSON(w, http.StatusOK, map[string]string{"value": value})
	})

	req := httptest.NewRequest("GET", "/test?key=test-value", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"success":true,"data":{"value":"test-value"}}`, w.Body.String())
}

func TestWrapHandlerWithParams(t *testing.T) {
	router := NewChiRouter()

	handler := func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
		return RespondWithJSON(w, http.StatusOK, params)
	}

	router.Get("/users/{id}/posts/{postId}", WrapHandlerWithParams(handler))

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"success":true,"data":{"id":"123","postId":"456"}}`, w.Body.String())
}

func TestWrapErrorHandler(t *testing.T) {
	router := NewChiRouter()

	handler := func(w http.ResponseWriter, r *http.Request) error {
		return RespondWithError(w, http.StatusBadRequest, "test error")
	}

	router.Get("/error", WrapErrorHandler(handler))

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"success":false,"error":"test error"}`, w.Body.String())
}

func TestRouteGroups(t *testing.T) {
	router := NewChiRouter()

	// Create an API group with a prefix
	router.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			// Add middleware to the group
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-API-Version", "v1")
					next.ServeHTTP(w, r)
				})
			})

			// Add routes to the group
			r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "get users"})
			})

			r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{"message": "create user"})
			})

			// Create a nested group
			r.Route("/admin", func(r chi.Router) {
				r.Use(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Admin", "true")
						next.ServeHTTP(w, r)
					})
				})

				r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"message": "admin stats"})
				})
			})
		})
	})

	// Test cases
	tests := []struct {
		name            string
		method          string
		path            string
		expectedCode    int
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{
			name:         "GET /api/users",
			method:       "GET",
			path:         "/api/users",
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"get users"}`,
			expectedHeaders: map[string]string{
				"X-API-Version": "v1",
			},
		},
		{
			name:         "POST /api/users",
			method:       "POST",
			path:         "/api/users",
			expectedCode: http.StatusCreated,
			expectedBody: `{"message":"create user"}`,
			expectedHeaders: map[string]string{
				"X-API-Version": "v1",
			},
		},
		{
			name:         "GET /api/admin/stats",
			method:       "GET",
			path:         "/api/admin/stats",
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"admin stats"}`,
			expectedHeaders: map[string]string{
				"X-API-Version": "v1",
				"X-Admin":       "true",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.Router().ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			body := strings.TrimSpace(w.Body.String())
			if body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}

			for key, value := range tt.expectedHeaders {
				if got := w.Header().Get(key); got != value {
					t.Errorf("expected header %s=%s, got %s", key, value, got)
				}
			}
		})
	}
}

func TestChiRouterErrorHandler(t *testing.T) {
	router := NewChiRouter()

	errorHandler := func(w http.ResponseWriter, r *http.Request) error {
		return io.EOF
	}

	router.Get("/error", WrapErrorHandler(errorHandler))

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp Response
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, io.EOF.Error(), resp.Error)
}

func TestChiRouterGroup(t *testing.T) {
	router := NewChiRouter()

	router.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]string{"message": "group"})
			})
		})
	})

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "group", resp["message"])
}

func TestJSONMiddleware(t *testing.T) {
	router := NewChiRouter()

	router.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/json", nil)
	w := httptest.NewRecorder()
	router.Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
