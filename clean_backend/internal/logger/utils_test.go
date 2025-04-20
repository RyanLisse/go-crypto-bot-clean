package logger

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestWithFields(t *testing.T) {
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

	// Create fields
	fields := Fields{
		"string_field": "string value",
		"int_field":    42,
		"bool_field":   true,
		"float_field":  3.14,
	}

	// Get a logger with fields
	fieldsLogger := WithFields(fields)

	// Log a message
	fieldsLogger.Info().Msg("test message")

	// Close the writer to flush the buffer
	w.Close()

	// Read the output from the pipe
	output := make([]byte, 1024)
	n, _ := r.Read(output)
	logOutput := string(output[:n])

	// Check that the message was logged with the fields
	if !strings.Contains(logOutput, "test message") {
		t.Errorf("Expected log to contain 'test message', got %q", logOutput)
	}
	if !strings.Contains(logOutput, "string_field") || !strings.Contains(logOutput, "string value") {
		t.Errorf("Expected log to contain string field, got %q", logOutput)
	}
	if !strings.Contains(logOutput, "int_field") || !strings.Contains(logOutput, "42") {
		t.Errorf("Expected log to contain int field, got %q", logOutput)
	}
	if !strings.Contains(logOutput, "bool_field") || !strings.Contains(logOutput, "true") {
		t.Errorf("Expected log to contain bool field, got %q", logOutput)
	}
	if !strings.Contains(logOutput, "float_field") || !strings.Contains(logOutput, "3.14") {
		t.Errorf("Expected log to contain float field, got %q", logOutput)
	}
}

func TestInfoWithFields(t *testing.T) {
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

	// Create fields
	fields := Fields{
		"string_field": "string value",
		"int_field":    42,
	}

	// Log a message with fields
	InfoWithFields("test message", fields)

	// Close the writer to flush the buffer
	w.Close()

	// Read the output from the pipe
	output := make([]byte, 1024)
	n, _ := r.Read(output)
	logOutput := string(output[:n])

	// Check that the message was logged with the fields
	if !strings.Contains(logOutput, "test message") {
		t.Errorf("Expected log to contain 'test message', got %q", logOutput)
	}
	if !strings.Contains(logOutput, "string_field") || !strings.Contains(logOutput, "string value") {
		t.Errorf("Expected log to contain string field, got %q", logOutput)
	}
	if !strings.Contains(logOutput, "int_field") || !strings.Contains(logOutput, "42") {
		t.Errorf("Expected log to contain int field, got %q", logOutput)
	}
}

func TestFunctionName(t *testing.T) {
	name := FunctionName()
	if !strings.Contains(name, "TestFunctionName") {
		t.Errorf("Expected function name to contain 'TestFunctionName', got %q", name)
	}
}

func TestWithContextNoValues(t *testing.T) {
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

	// Create a context with no values
	ctx := context.Background()

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

	// Check that the message was logged without context values
	if !strings.Contains(logOutput, "test message") {
		t.Errorf("Expected log to contain 'test message', got %q", logOutput)
	}
	if strings.Contains(logOutput, "request_id") {
		t.Errorf("Expected log to not contain 'request_id', got %q", logOutput)
	}
	if strings.Contains(logOutput, "user_id") {
		t.Errorf("Expected log to not contain 'user_id', got %q", logOutput)
	}
}
