package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/rs/zerolog"
)

func TestNewRouter_HealthAndRootEndpoints(t *testing.T) {
	logger := zerolog.Nop()
	cfg := &config.Config{Version: "test-version"}
	r := NewRouter(cfg, &logger)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Health endpoint
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("/health request error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("/health wrong status: got %d", resp.StatusCode)
	}
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, resp.Body)
	if !strings.Contains(buf.String(), "test-version") {
		t.Errorf("/health missing version")
	}

	// Root-test endpoint
	resp2, err := http.Get(ts.URL + "/root-test")
	if err != nil {
		t.Fatalf("/root-test request error: %v", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("/root-test wrong status: got %d", resp2.StatusCode)
	}
}
