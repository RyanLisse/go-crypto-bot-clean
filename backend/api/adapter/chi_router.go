package adapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ChiRouter wraps chi.Router to provide additional functionality
type ChiRouter struct {
	router chi.Router
}

// NewChiRouter creates a new ChiRouter instance with default middleware
func NewChiRouter() *ChiRouter {
	r := chi.NewRouter()

	// Add default middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(JSONMiddleware)

	return &ChiRouter{router: r}
}

// Router returns the underlying chi.Router instance
func (cr *ChiRouter) Router() chi.Router {
	return cr.router
}

// Use adds middleware to the router
func (cr *ChiRouter) Use(middlewares ...func(http.Handler) http.Handler) {
	cr.router.Use(middlewares...)
}

// Group creates a new route group
func (cr *ChiRouter) Group(fn func(r chi.Router)) chi.Router {
	return cr.router.Group(fn)
}

// Handle registers a route with a specific HTTP method
func (cr *ChiRouter) Handle(pattern string, handler http.Handler) {
	cr.router.Handle(pattern, handler)
}

// Method registers a route with a specific HTTP method
func (cr *ChiRouter) Method(method, pattern string, handler http.HandlerFunc) {
	cr.router.Method(method, pattern, handler)
}

// Get registers a GET route
func (cr *ChiRouter) Get(pattern string, handler http.HandlerFunc) {
	cr.router.Get(pattern, handler)
}

// Post registers a POST route
func (cr *ChiRouter) Post(pattern string, handler http.HandlerFunc) {
	cr.router.Post(pattern, handler)
}

// Put registers a PUT route
func (cr *ChiRouter) Put(pattern string, handler http.HandlerFunc) {
	cr.router.Put(pattern, handler)
}

// Delete registers a DELETE route
func (cr *ChiRouter) Delete(pattern string, handler http.HandlerFunc) {
	cr.router.Delete(pattern, handler)
}

// WithContext adds context values to requests handled by the router
func (cr *ChiRouter) WithContext(key interface{}, value interface{}) *ChiRouter {
	cr.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	return cr
}

// JSONMiddleware sets the content type header to application/json
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// RespondWithJSON sends a JSON response with the given status code and data
func RespondWithJSON(w http.ResponseWriter, code int, data interface{}) error {
	response := Response{
		Success: code >= 200 && code < 300,
		Data:    data,
	}
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

// RespondWithError sends an error response with the given status code and message
func RespondWithError(w http.ResponseWriter, code int, message string) error {
	response := Response{
		Success: false,
		Error:   message,
	}
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(response)
}

// GetParam retrieves a URL parameter from the request context
func GetParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// GetQuery retrieves a query parameter from the request URL
func GetQuery(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetBody reads and returns the request body as a byte slice
func GetBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

// ErrorHandler wraps a handler function that returns an error
type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

// WrapErrorHandler converts an ErrorHandler to http.HandlerFunc with error handling
func WrapErrorHandler(handler ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

// HandlerWithParams is a function type for handling URL parameters
type HandlerWithParams func(w http.ResponseWriter, r *http.Request, params map[string]string) error

// WrapHandlerWithParams converts a HandlerWithParams to http.HandlerFunc
func WrapHandlerWithParams(handler HandlerWithParams) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := make(map[string]string)

		// Add route parameters from Chi context
		rctx := chi.RouteContext(r.Context())
		if rctx != nil {
			for i, key := range rctx.URLParams.Keys {
				params[key] = rctx.URLParams.Values[i]
			}
		}

		if err := handler(w, r, params); err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
}
