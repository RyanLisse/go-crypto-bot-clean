package logger

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger_NotNil(t *testing.T) {
	l := NewLogger()
	if l == nil {
		t.Error("NewLogger() returned nil")
	}
}

func TestNewLoggerWithLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{"debug", "debug", zerolog.DebugLevel},
		{"info", "info", zerolog.InfoLevel},
		{"warn", "warn", zerolog.WarnLevel},
		{"error", "error", zerolog.ErrorLevel},
		{"fatal", "fatal", zerolog.FatalLevel},
		{"panic", "panic", zerolog.PanicLevel},
		{"invalid", "invalid", zerolog.InfoLevel}, // Default to info for invalid levels
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewLoggerWithLevel(tt.level)
			if zerolog.GlobalLevel() != tt.want {
				t.Errorf("NewLoggerWithLevel(%q) set global level to %v, want %v", tt.level, zerolog.GlobalLevel(), tt.want)
			}
		})
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	// Create a config that writes to stdout
	config := DefaultConfig()
	config.Format = "json"   // Use JSON format for easier testing
	config.Output = "stdout" // We'll replace os.Stdout with our pipe

	// Save the original stdout and restore it after the test
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a logger with the config
	logger := NewLoggerWithConfig(config)

	// Log a message
	logger.Info().Msg("test message")

	// Close the writer to flush the buffer
	w.Close()

	// Read the output from the pipe
	output := make([]byte, 1024)
	n, _ := r.Read(output)
	logOutput := string(output[:n])

	// Check that the message was logged
	if !strings.Contains(logOutput, "test message") {
		t.Errorf("Expected log to contain 'test message', got %q", logOutput)
	}

	// Check that the log is in JSON format
	if !strings.Contains(logOutput, "{") || !strings.Contains(logOutput, "}") {
		t.Errorf("Expected log to be in JSON format, got %q", logOutput)
	}
}

func TestWithContext(t *testing.T) {
	// Create a config that writes to stdout
	config := DefaultConfig()
	config.Format = "json"   // Use JSON format for easier testing
	config.Output = "stdout" // We'll replace os.Stdout with our pipe

	// Save the original stdout and restore it after the test
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a logger with the config and set it as the global logger
	logger := NewLoggerWithConfig(config)
	SetGlobalLogger(logger)

	// Create a context with request ID and user ID
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "test-request-id")
	ctx = context.WithValue(ctx, "user_id", "test-user-id")

	// Get a logger with the context
	contextLogger := WithContext(ctx)

	// Log a message
	contextLogger.Info().Msg("test message")

	// Close the writer to flush the buffer
	w.Close()

	// Read the output from the pipe
	output := make([]byte, 1024)
	n, _ := r.Read(output)
	logOutput := string(output[:n])

	// Check that the message was logged with the context values
	if !strings.Contains(logOutput, "test message") {
		t.Errorf("Expected log to contain 'test message', got %q", logOutput)
	}
	if !strings.Contains(logOutput, "test-request-id") {
		t.Errorf("Expected log to contain 'test-request-id', got %q", logOutput)
	}
	if !strings.Contains(logOutput, "test-user-id") {
		t.Errorf("Expected log to contain 'test-user-id', got %q", logOutput)
	}
}

func TestWithComponent(t *testing.T) {
	// Create a config that writes to stdout
	config := DefaultConfig()
	config.Format = "json"   // Use JSON format for easier testing
	config.Output = "stdout" // We'll replace os.Stdout with our pipe

	// Save the original stdout and restore it after the test
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a logger with the config and set it as the global logger
	logger := NewLoggerWithConfig(config)
	SetGlobalLogger(logger)

	// Get a logger with a component
	componentLogger := WithComponent("test-component")

	// Log a message
	componentLogger.Info().Msg("test message")

	// Close the writer to flush the buffer
	w.Close()

	// Read the output from the pipe
	output := make([]byte, 1024)
	n, _ := r.Read(output)
	logOutput := string(output[:n])

	// Check that the message was logged with the component
	if !strings.Contains(logOutput, "test message") {
		t.Errorf("Expected log to contain 'test message', got %q", logOutput)
	}
	if !strings.Contains(logOutput, "test-component") {
		t.Errorf("Expected log to contain 'test-component', got %q", logOutput)
	}
}

func TestGlobalLoggerFunctions(t *testing.T) {
	// Create a config that writes to stdout
	config := DefaultConfig()
	config.Format = "json"   // Use JSON format for easier testing
	config.Output = "stdout" // We'll replace os.Stdout with our pipe

	// Save the original stdout and restore it after the test
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a logger with the config and set it as the global logger
	logger := NewLoggerWithConfig(config)
	SetGlobalLogger(logger)

	// Test the global logger functions
	Info("info message")
	Warn("warn message")
	Debug("debug message") // This won't be logged at the default info level

	// Close the writer to flush the buffer
	w.Close()

	// Read the output from the pipe
	output := make([]byte, 1024)
	n, _ := r.Read(output)
	logOutput := string(output[:n])

	// Check that the messages were logged
	if !strings.Contains(logOutput, "info message") {
		t.Errorf("Expected log to contain 'info message', got %q", logOutput)
	}
	if !strings.Contains(logOutput, "warn message") {
		t.Errorf("Expected log to contain 'warn message', got %q", logOutput)
	}

	// Debug message should not be logged at info level
	if strings.Contains(logOutput, "debug message") {
		t.Errorf("Expected log to not contain 'debug message', got %q", logOutput)
	}
}
