package logger

import (
	"testing"
)

func TestNewLogger_NotNil(t *testing.T) {
	l := NewLogger()
	if l == nil {
		t.Error("NewLogger() returned nil")
	}
}
