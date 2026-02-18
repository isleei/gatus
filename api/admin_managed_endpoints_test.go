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

func TestManagedEndpointsCRUD(t *testing.T) {
	cfg := loadManagedEndpointsTestConfig(t)
	router := New(cfg).Router()

	getCode, getBody := runManagedEndpointRequest(t, router, http.MethodGet, "/api/v1/admin/endpoints", "")
	if getCode != http.StatusOK {
		t.Fatalf("expected GET code %d, got %d", http.StatusOK, getCode)
	}
	if !strings.Contains(getBody, `"key":"core_frontend"`) {
		t.Fatalf("expected body to include core_frontend key, got: %s", getBody)
	}

	createPayload := `{
  "enabled": true,
  "name": "backend",
  "group": "core",
  "url": "https://example.org/health",
  "method": "GET",
  "interval": "45s",
  "conditions": ["[STATUS] == 200"]
}`
	createCode, createBody := runManagedEndpointRequest(t, router, http.MethodPost, "/api/v1/admin/endpoints", createPayload)
	if createCode != http.StatusCreated {
		t.Fatalf("expected POST code %d, got %d (%s)", http.StatusCreated, createCode, createBody)
	}
	if !strings.Contains(createBody, `"key":"core_backend"`) {
		t.Fatalf("expected response to include core_backend key, got: %s", createBody)
	}

	overlayAfterCreate, err := config.ReadManagedOverlay(cfg.LoadedConfigPath())
	if err != nil {
		t.Fatalf("unexpected error reading overlay: %v", err)
	}
	if overlayAfterCreate == nil || len(overlayAfterCreate.Endpoints) != 2 {
		t.Fatalf("expected overlay to contain 2 endpoints after create, got %+v", overlayAfterCreate)
	}

	updatePayload := `{
  "enabled": true,
  "name": "api",
  "group": "core",
  "url": "https://example.org/v2/health",
  "method": "POST",
  "interval": "1m",
  "conditions": ["[STATUS] == 200", "[RESPONSE_TIME] < 500"],
  "headers": {"X-Test": "1"}
}`
	updateCode, updateBody := runManagedEndpointRequest(t, router, http.MethodPut, "/api/v1/admin/endpoints/core_backend", updatePayload)
	if updateCode != http.StatusOK {
		t.Fatalf("expected PUT code %d, got %d (%s)", http.StatusOK, updateCode, updateBody)
	}
	if !strings.Contains(updateBody, `"key":"core_api"`) {
		t.Fatalf("expected response to include renamed key core_api, got: %s", updateBody)
	}

	deleteCode, deleteBody := runManagedEndpointRequest(t, router, http.MethodDelete, "/api/v1/admin/endpoints/core_api", "")
	if deleteCode != http.StatusNoContent {
		t.Fatalf("expected DELETE code %d, got %d (%s)", http.StatusNoContent, deleteCode, deleteBody)
	}

	overlayAfterDelete, err := config.ReadManagedOverlay(cfg.LoadedConfigPath())
	if err != nil {
		t.Fatalf("unexpected error reading overlay after delete: %v", err)
	}
	if overlayAfterDelete == nil || len(overlayAfterDelete.Endpoints) != 1 {
		t.Fatalf("expected overlay to contain 1 endpoint after delete, got %+v", overlayAfterDelete)
	}
	if overlayAfterDelete.Endpoints[0].Key() != "core_frontend" {
		t.Fatalf("expected remaining endpoint to be core_frontend, got %s", overlayAfterDelete.Endpoints[0].Key())
	}
}

func TestManagedEndpointsCreateInvalidInterval(t *testing.T) {
	cfg := loadManagedEndpointsTestConfig(t)
	router := New(cfg).Router()
	payload := `{
  "name": "backend",
  "group": "core",
  "url": "https://example.org/health",
  "interval": "not-a-duration",
  "conditions": ["[STATUS] == 200"]
}`
	code, body := runManagedEndpointRequest(t, router, http.MethodPost, "/api/v1/admin/endpoints", payload)
	if code != http.StatusBadRequest {
		t.Fatalf("expected POST code %d, got %d (%s)", http.StatusBadRequest, code, body)
	}
	if !strings.Contains(body, "invalid interval") {
		t.Fatalf("expected invalid interval error, got: %s", body)
	}
}

func runManagedEndpointRequest(t *testing.T, router *fiber.App, method, path, body string) (int, string) {
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

func loadManagedEndpointsTestConfig(t *testing.T) *config.Config {
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
