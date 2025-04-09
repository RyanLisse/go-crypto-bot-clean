package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global logger instance
var Logger *LoggerWrapper

// LoggerWrapper wraps a zap.Logger with additional functionality
type LoggerWrapper struct {
	Logger *zap.Logger
}

// init initializes the global logger
func init() {
	// Create a production logger
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Create the wrapper
	Logger = &LoggerWrapper{
		Logger: logger,
	}
}

// With returns a logger with the given fields
func (l *LoggerWrapper) With(fields ...zap.Field) *LoggerWrapper {
	return &LoggerWrapper{
		Logger: l.Logger.With(fields...),
	}
}

// Debug logs a debug message
func (l *LoggerWrapper) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Info logs an info message
func (l *LoggerWrapper) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *LoggerWrapper) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Error logs an error message
func (l *LoggerWrapper) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *LoggerWrapper) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
}
