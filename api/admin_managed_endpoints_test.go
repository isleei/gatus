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
  "conditions": ["[STATUS] == 200"],
  "client": {
    "timeout": "15s",
    "insecure": true,
    "ignoreRedirect": true
  },
  "certificate": {
    "enabled": true,
    "expirationThreshold": "72h"
  },
  "tamper": {
    "enabled": true,
    "baselineSamples": 5,
    "driftThresholdPercent": 15,
    "consecutiveBreaches": 2,
    "requiredSubstrings": ["<title>Status</title>", "在线状态页"],
    "forbiddenSubstrings": ["博彩", "hack"]
  }
}`
	createCode, createBody := runManagedEndpointRequest(t, router, http.MethodPost, "/api/v1/admin/endpoints", createPayload)
	if createCode != http.StatusCreated {
		t.Fatalf("expected POST code %d, got %d (%s)", http.StatusCreated, createCode, createBody)
	}
	if !strings.Contains(createBody, `"key":"core_backend"`) {
		t.Fatalf("expected response to include core_backend key, got: %s", createBody)
	}
	if !strings.Contains(createBody, `"expirationThreshold":"72h"`) {
		t.Fatalf("expected response to include certificate expiration threshold, got: %s", createBody)
	}
	if !strings.Contains(createBody, `"driftThresholdPercent":15`) {
		t.Fatalf("expected response to include tamper configuration, got: %s", createBody)
	}

	overlayAfterCreate, err := config.ReadManagedOverlay(cfg.LoadedConfigPath())
	if err != nil {
		t.Fatalf("unexpected error reading overlay: %v", err)
	}
	if overlayAfterCreate == nil || len(overlayAfterCreate.Endpoints) != 2 {
		t.Fatalf("expected overlay to contain 2 endpoints after create, got %+v", overlayAfterCreate)
	}
	var createdEndpointFound bool
	for _, ep := range overlayAfterCreate.Endpoints {
		if ep.Key() != "core_backend" {
			continue
		}
		createdEndpointFound = true
		if ep.ClientConfig == nil || ep.ClientConfig.Timeout.String() != "15s" || !ep.ClientConfig.Insecure || !ep.ClientConfig.IgnoreRedirect {
			t.Fatalf("expected client config to be persisted, got %+v", ep.ClientConfig)
		}
		if ep.TamperConfig == nil || !ep.TamperConfig.Enabled || ep.TamperConfig.BaselineSamples != 5 || ep.TamperConfig.DriftThresholdPercent != 15 || ep.TamperConfig.ConsecutiveBreaches != 2 {
			t.Fatalf("expected tamper config to be persisted, got %+v", ep.TamperConfig)
		}
		if len(ep.TamperConfig.RequiredSubstrings) != 2 || len(ep.TamperConfig.ForbiddenSubstrings) != 2 {
			t.Fatalf("expected tamper required/forbidden substrings to be persisted, got %+v", ep.TamperConfig)
		}
		joinedConditions := make([]string, 0, len(ep.Conditions))
		for _, condition := range ep.Conditions {
			joinedConditions = append(joinedConditions, string(condition))
		}
		if !strings.Contains(strings.Join(joinedConditions, "\n"), "[CERTIFICATE_EXPIRATION] > 72h") {
			t.Fatalf("expected generated certificate condition, got %+v", joinedConditions)
		}
	}
	if !createdEndpointFound {
		t.Fatal("expected created endpoint core_backend in overlay")
	}

	updatePayload := `{
  "enabled": true,
  "name": "api",
  "group": "core",
  "url": "https://example.org/v2/health",
  "method": "POST",
  "interval": "1m",
  "conditions": ["[STATUS] == 200", "[RESPONSE_TIME] < 500"],
  "headers": {"X-Test": "1"},
  "client": {
    "timeout": "20s",
    "insecure": false,
    "ignoreRedirect": false
  },
  "certificate": {
    "enabled": true,
    "expirationThreshold": "24h"
  },
  "tamper": {
    "enabled": true,
    "baselineSamples": 10,
    "driftThresholdPercent": 20,
    "consecutiveBreaches": 3,
    "requiredSubstrings": ["<title>API</title>"],
    "forbiddenSubstrings": ["博彩"]
  },
  "alerts": [
    {
      "type": "slack",
      "minimumReminderInterval": "1h"
    }
  ],
  "ui": {
    "resolveSuccessfulConditions": true
  }
}`
	updateCode, updateBody := runManagedEndpointRequest(t, router, http.MethodPut, "/api/v1/admin/endpoints/core_backend", updatePayload)
	if updateCode != http.StatusOK {
		t.Fatalf("expected PUT code %d, got %d (%s)", http.StatusOK, updateCode, updateBody)
	}
	if !strings.Contains(updateBody, `"key":"core_api"`) {
		t.Fatalf("expected response to include renamed key core_api, got: %s", updateBody)
	}
	if !strings.Contains(updateBody, `"minimumReminderInterval":"1h0m0s"`) {
		t.Fatalf("expected response to include normalized reminder interval, got: %s", updateBody)
	}
	if !strings.Contains(updateBody, `"resolveSuccessfulConditions":true`) {
		t.Fatalf("expected response to include endpoint ui config, got: %s", updateBody)
	}
	if !strings.Contains(updateBody, `"timeout":"20s"`) {
		t.Fatalf("expected response to include client timeout, got: %s", updateBody)
	}
	if !strings.Contains(updateBody, `"expirationThreshold":"24h"`) {
		t.Fatalf("expected response to include updated certificate expiration threshold, got: %s", updateBody)
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
