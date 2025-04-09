package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	"go-crypto-bot-clean/backend/internal/auth"
	"go-crypto-bot-clean/backend/internal/logging"

	"go.uber.org/zap"
)

// ServiceMode represents the current mode of the service
type ServiceMode int32

const (
	// NormalMode means the service is operating normally
	NormalMode ServiceMode = iota
	// ReadOnlyMode means the service is in read-only mode
	ReadOnlyMode
	// MaintenanceMode means the service is in maintenance mode
	MaintenanceMode
)

// GracefulDegradationConfig holds configuration for graceful degradation
type GracefulDegradationConfig struct {
	// ReadOnlyThreshold is the threshold for switching to read-only mode
	ReadOnlyThreshold int
	// MaintenanceThreshold is the threshold for switching to maintenance mode
	MaintenanceThreshold int
	// RecoveryInterval is the interval for recovery checks
	RecoveryInterval time.Duration
	// Logger is the logger to use
	Logger *zap.Logger
	// ReadOnlyPaths are paths that are always allowed in read-only mode
	ReadOnlyPaths []string
	// CriticalPaths are paths that are always allowed in any mode
	CriticalPaths []string
}

// DefaultGracefulDegradationConfig returns the default graceful degradation configuration
func DefaultGracefulDegradationConfig() GracefulDegradationConfig {
	return GracefulDegradationConfig{
		ReadOnlyThreshold:    10,
		MaintenanceThreshold: 20,
		RecoveryInterval:     30 * time.Second,
		Logger:               nil,
		ReadOnlyPaths:        []string{"/api/v1/health", "/api/v1/status", "/api/v1/metrics"},
		CriticalPaths:        []string{"/api/v1/health"},
	}
}

// GracefulDegradation implements graceful degradation
type GracefulDegradation struct {
	config           GracefulDegradationConfig
	mode             int32
	errorCount       int32
	lastRecoveryTime time.Time
	logger           *zap.Logger
}

// NewGracefulDegradation creates a new graceful degradation instance
func NewGracefulDegradation(config GracefulDegradationConfig) *GracefulDegradation {
	if config.Logger == nil {
		config.Logger = logging.Logger.Logger
	}

	return &GracefulDegradation{
		config:           config,
		mode:             int32(NormalMode),
		errorCount:       0,
		lastRecoveryTime: time.Now(),
		logger:           config.Logger,
	}
}

// GracefulDegradationMiddleware adds graceful degradation to a handler
func GracefulDegradationMiddleware(gd *GracefulDegradation) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if recovery is needed
			gd.checkRecovery()

			// Get current mode
			mode := ServiceMode(atomic.LoadInt32(&gd.mode))

			// Check if request is allowed in current mode
			if !gd.isRequestAllowed(r, mode) {
				// Request is not allowed in current mode
				var errResp *auth.AuthError
				switch mode {
				case ReadOnlyMode:
					errResp = auth.NewAuthError(
						auth.ErrorTypeReadOnly,
						"Service is in read-only mode",
						http.StatusServiceUnavailable,
					).WithDetails(map[string]interface{}{
						"service_mode": "read_only",
						"allowed_methods": []string{"GET", "HEAD", "OPTIONS"},
					}).WithHelp("Only read operations are allowed at this time").
						WithRequestID(GetRequestID(r.Context()))
				case MaintenanceMode:
					errResp = auth.NewAuthError(
						auth.ErrorTypeMaintenance,
						"Service is in maintenance mode",
						http.StatusServiceUnavailable,
					).WithDetails(map[string]interface{}{
						"service_mode": "maintenance",
					}).WithHelp("The service is temporarily unavailable for maintenance").
						WithRequestID(GetRequestID(r.Context()))
				}

				// Set response headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)

				// Write error response
				errResp.WriteJSON(w)
				return
			}

			// Handle the request
			next.ServeHTTP(w, r)
		})
	}
}

// RecordError records an error and potentially degrades service
func (gd *GracefulDegradation) RecordError() {
	// Increment error count
	count := atomic.AddInt32(&gd.errorCount, 1)

	// Check if we need to degrade
	currentMode := ServiceMode(atomic.LoadInt32(&gd.mode))
	if currentMode == NormalMode && count >= int32(gd.config.ReadOnlyThreshold) {
		// Switch to read-only mode
		atomic.StoreInt32(&gd.mode, int32(ReadOnlyMode))
		gd.logger.Warn("Service degraded to read-only mode",
			zap.Int32("error_count", count),
			zap.Int("threshold", gd.config.ReadOnlyThreshold),
		)
	} else if currentMode == ReadOnlyMode && count >= int32(gd.config.MaintenanceThreshold) {
		// Switch to maintenance mode
		atomic.StoreInt32(&gd.mode, int32(MaintenanceMode))
		gd.logger.Warn("Service degraded to maintenance mode",
			zap.Int32("error_count", count),
			zap.Int("threshold", gd.config.MaintenanceThreshold),
		)
	}
}

// checkRecovery checks if recovery is needed
func (gd *GracefulDegradation) checkRecovery() {
	// Check if recovery interval has passed
	if time.Since(gd.lastRecoveryTime) < gd.config.RecoveryInterval {
		return
	}

	// Reset recovery time
	gd.lastRecoveryTime = time.Now()

	// Get current mode
	currentMode := ServiceMode(atomic.LoadInt32(&gd.mode))
	if currentMode == NormalMode {
		// Already in normal mode, just reset error count
		atomic.StoreInt32(&gd.errorCount, 0)
		return
	}

	// Try to recover one level at a time
	if currentMode == MaintenanceMode {
		// Try to recover to read-only mode
		atomic.StoreInt32(&gd.mode, int32(ReadOnlyMode))
		atomic.StoreInt32(&gd.errorCount, int32(gd.config.ReadOnlyThreshold))
		gd.logger.Info("Service recovered from maintenance to read-only mode")
	} else if currentMode == ReadOnlyMode {
		// Try to recover to normal mode
		atomic.StoreInt32(&gd.mode, int32(NormalMode))
		atomic.StoreInt32(&gd.errorCount, 0)
		gd.logger.Info("Service recovered from read-only to normal mode")
	}
}

// isRequestAllowed checks if a request is allowed in the current mode
func (gd *GracefulDegradation) isRequestAllowed(r *http.Request, mode ServiceMode) bool {
	// Check if path is in critical paths (always allowed)
	path := r.URL.Path
	for _, criticalPath := range gd.config.CriticalPaths {
		if path == criticalPath {
			return true
		}
	}

	// Check mode-specific rules
	switch mode {
	case NormalMode:
		// All requests allowed in normal mode
		return true
	case ReadOnlyMode:
		// Only GET, HEAD, OPTIONS requests allowed in read-only mode
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			return true
		}
		// Check if path is in read-only paths
		for _, readOnlyPath := range gd.config.ReadOnlyPaths {
			if path == readOnlyPath {
				return true
			}
		}
		return false
	case MaintenanceMode:
		// Only critical paths allowed in maintenance mode
		return false
	default:
		return false
	}
}

// GetMode returns the current service mode
func (gd *GracefulDegradation) GetMode() ServiceMode {
	return ServiceMode(atomic.LoadInt32(&gd.mode))
}

// SetMode sets the service mode
func (gd *GracefulDegradation) SetMode(mode ServiceMode) {
	atomic.StoreInt32(&gd.mode, int32(mode))
	gd.logger.Info("Service mode manually set",
		zap.Int32("mode", int32(mode)),
	)
}

// GetErrorCount returns the current error count
func (gd *GracefulDegradation) GetErrorCount() int32 {
	return atomic.LoadInt32(&gd.errorCount)
}

// ResetErrorCount resets the error count
func (gd *GracefulDegradation) ResetErrorCount() {
	atomic.StoreInt32(&gd.errorCount, 0)
	gd.logger.Info("Error count manually reset")
}
