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
	"github.com/gofiber/fiber/v2"
)

func TestManagedConfigCRUD(t *testing.T) {
	cfg := loadManagedConfigTestConfig(t)
	router := New(cfg).Router()

	// GET initial state
	getCode, getBody := runManagedConfigRequest(t, router, http.MethodGet, "/api/v1/admin/managed-config", "")
	if getCode != http.StatusOK {
		t.Fatalf("expected GET code %d, got %d", http.StatusOK, getCode)
	}
	if !strings.Contains(getBody, `"overlayPath"`) {
		t.Fatalf("expected body to include overlayPath, got: %s", getBody)
	}

	// PUT valid payload
	putPayload := `{
  "endpoints": [
    {
      "Name": "api",
      "Group": "core",
      "URL": "https://example.org/api/health",
      "Interval": 30000000000,
      "Conditions": ["[STATUS] == 200"]
    }
  ],
  "externalEndpoints": [],
  "suites": []
}`
	putCode, putBody := runManagedConfigRequest(t, router, http.MethodPut, "/api/v1/admin/managed-config", putPayload)
	if putCode != http.StatusOK {
		t.Fatalf("expected PUT code %d, got %d (%s)", http.StatusOK, putCode, putBody)
	}
	if !strings.Contains(putBody, "Managed configuration saved") {
		t.Fatalf("expected success message, got: %s", putBody)
	}

	// GET after PUT to verify
	getCode2, getBody2 := runManagedConfigRequest(t, router, http.MethodGet, "/api/v1/admin/managed-config", "")
	if getCode2 != http.StatusOK {
		t.Fatalf("expected GET code %d after PUT, got %d", http.StatusOK, getCode2)
	}
	if !strings.Contains(getBody2, `"Name":"api"`) {
		t.Fatalf("expected body to include api endpoint after PUT, got: %s", getBody2)
	}

	// DELETE
	deleteCode, _ := runManagedConfigRequest(t, router, http.MethodDelete, "/api/v1/admin/managed-config", "")
	if deleteCode != http.StatusNoContent {
		t.Fatalf("expected DELETE code %d, got %d", http.StatusNoContent, deleteCode)
	}

	// GET after DELETE to verify overlay is gone
	getCode3, getBody3 := runManagedConfigRequest(t, router, http.MethodGet, "/api/v1/admin/managed-config", "")
	if getCode3 != http.StatusOK {
		t.Fatalf("expected GET code %d after DELETE, got %d", http.StatusOK, getCode3)
	}
	// After delete, it should fall back to the base YAML config endpoint
	if !strings.Contains(getBody3, `"Name":"frontend"`) {
		t.Fatalf("expected body to contain base YAML endpoint after delete, got: %s", getBody3)
	}
}

func TestManagedConfigPutInvalidPayload(t *testing.T) {
	cfg := loadManagedConfigTestConfig(t)
	router := New(cfg).Router()

	// Empty endpoints and suites should fail validation
	badPayload := `{"endpoints":[],"externalEndpoints":[],"suites":[]}`
	code, body := runManagedConfigRequest(t, router, http.MethodPut, "/api/v1/admin/managed-config", badPayload)
	if code != http.StatusBadRequest {
		t.Fatalf("expected code %d for invalid payload, got %d (%s)", http.StatusBadRequest, code, body)
	}
}

func runManagedConfigRequest(t *testing.T, router *fiber.App, method, path, body string) (int, string) {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if len(body) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := router.Test(request)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("unexpected body read error: %v", err)
	}
	return response.StatusCode, string(responseBody)
}

func loadManagedConfigTestConfig(t *testing.T) *config.Config {
	t.Helper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	configYAML := `endpoints:
  - name: frontend
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
