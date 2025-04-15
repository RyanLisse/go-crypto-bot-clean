//go:build !turso

// Package turso provides a stub implementation when Turso is not enabled
package turso

import (
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog"
)

// TursoDB is a stub implementation when Turso is not enabled
type TursoDB struct {
	db     *sql.DB
	logger *zerolog.Logger
}

// ErrTursoNotEnabled is returned when Turso functionality is used but not enabled
var ErrTursoNotEnabled = errors.New("turso is not enabled in this build")

// NewTursoDB returns an error since Turso is not enabled
func NewTursoDB(primaryURL, authToken string, syncInterval time.Duration, logger *zerolog.Logger) (*TursoDB, error) {
	logger.Warn().Msg("Turso support is not enabled in this build")
	return nil, ErrTursoNotEnabled
}

// DB returns nil since Turso is not enabled
func (t *TursoDB) DB() *sql.DB {
	if t == nil || t.db == nil {
		return nil
	}
	return t.db
}

// Sync returns an error since Turso is not enabled
func (t *TursoDB) Sync() error {
	if t == nil || t.logger == nil {
		return ErrTursoNotEnabled
	}
	t.logger.Warn().Msg("Turso sync called but Turso is not enabled in this build")
	return ErrTursoNotEnabled
}

// Close is a no-op since Turso is not enabled
func (t *TursoDB) Close() error {
	if t == nil || t.logger == nil {
		return nil
	}
	t.logger.Warn().Msg("Turso close called but Turso is not enabled in this build")
	return nil
}
