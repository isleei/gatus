package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TwiN/gatus/v5/config"
)

func TestTriggerImmediateConfigurationReload(t *testing.T) {
	// Drain any pending reload signal
	select {
	case <-config.ImmediateReloadRequests():
	default:
	}

	cfg := loadReloadTestConfig(t)
	router := New(cfg).Router()

	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/reload", nil)
	response, err := router.Test(request)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d (%s)", http.StatusAccepted, response.StatusCode, string(body))
	}
	if !strings.Contains(string(body), "reload requested") {
		t.Fatalf("expected body to mention reload, got: %s", string(body))
	}

	// Verify the channel received a signal
	select {
	case <-config.ImmediateReloadRequests():
		// success
	default:
		t.Fatal("expected a signal on the reload channel after POST")
	}
}

func loadReloadTestConfig(t *testing.T) *config.Config {
	t.Helper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	configYAML := `endpoints:
  - name: test
    group: core
    url: https://example.org/health
    interval: 30s
    conditions:
      - "[STATUS] == 200"
`
	if err := os.WriteFile(configPath, []byte(configYAML), 0o600); err != nil {
		t.Fatalf("unexpected config write error: %v", err)
	}
	cfg, err := config.LoadConfiguration(configPath)
	if err != nil {
		t.Fatalf("unexpected config load error: %v", err)
	}
	return cfg
}
