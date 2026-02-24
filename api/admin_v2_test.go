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

func TestAdminV2SuitesAndExternalEndpointsCRUD(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	router := New(cfg).Router()

	createSuitePayload := `{
  "name": "checkout-flow",
  "group": "core",
  "interval": "2m",
  "timeout": "1m",
  "endpoints": [
    {
      "name": "step-1",
      "url": "https://example.org/step-1",
      "method": "GET",
      "conditions": ["[STATUS] == 200"],
      "alwaysRun": true
    }
  ]
}`
	code, body := runAdminV2Request(t, router, http.MethodPost, "/api/v1/admin/suites", createSuitePayload)
	if code != http.StatusCreated {
		t.Fatalf("expected suite create status %d, got %d (%s)", http.StatusCreated, code, body)
	}
	if !strings.Contains(body, `"key":"core_checkout-flow"`) {
		t.Fatalf("expected suite response to include key core_checkout-flow, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/suites", "")
	if code != http.StatusOK {
		t.Fatalf("expected suite list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"core_checkout-flow"`) {
		t.Fatalf("expected suite list to include created suite, got: %s", body)
	}

	createExternalPayload := `{
  "name": "worker-heartbeat",
  "group": "jobs",
  "token": "secret-token",
  "heartbeatInterval": "30s"
}`
	code, body = runAdminV2Request(t, router, http.MethodPost, "/api/v1/admin/external-endpoints", createExternalPayload)
	if code != http.StatusCreated {
		t.Fatalf("expected external endpoint create status %d, got %d (%s)", http.StatusCreated, code, body)
	}
	if !strings.Contains(body, `"key":"jobs_worker-heartbeat"`) {
		t.Fatalf("expected external endpoint response to include key jobs_worker-heartbeat, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/external-endpoints", "")
	if code != http.StatusOK {
		t.Fatalf("expected external endpoint list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"jobs_worker-heartbeat"`) {
		t.Fatalf("expected external endpoint list to include created endpoint, got: %s", body)
	}
}

func TestAdminV2BatchImportExportAndAudit(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	router := New(cfg).Router()

	code, body := runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/export?entityType=endpoint", "")
	if code != http.StatusOK {
		t.Fatalf("expected export status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"Name":"frontend"`) || !strings.Contains(body, `"Group":"core"`) {
		t.Fatalf("expected export to include frontend/core endpoint, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/monitors?entityType=all&group=all&enabled=all&status=all&sortBy=updatedAt&sortDir=desc&page=1&pageSize=50", "")
	if code != http.StatusOK {
		t.Fatalf("expected monitors status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"core_frontend"`) {
		t.Fatalf("expected monitor response to include core_frontend, got: %s", body)
	}

	batchDryRunPayload := `{
  "entityType": "endpoint",
  "keys": ["core_frontend"],
  "action": "set-group",
  "payload": {"group": "ops"},
  "dryRun": true
}`
	code, body = runAdminV2Request(t, router, http.MethodPost, "/api/v1/admin/monitors/batch", batchDryRunPayload)
	if code != http.StatusOK {
		t.Fatalf("expected batch dry run status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"dryRun":true`) {
		t.Fatalf("expected batch response to include dryRun=true, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/endpoints", "")
	if code != http.StatusOK {
		t.Fatalf("expected endpoints list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if strings.Contains(body, `"key":"ops_frontend"`) {
		t.Fatalf("expected dry run not to persist overlay changes, got: %s", body)
	}

	importDryRunPayload := `{
  "entityType": "endpoint",
  "mode": "merge",
  "dryRun": true,
  "data": {
    "endpoints": [
      {
        "name": "backend",
        "group": "core",
        "url": "https://example.org/backend",
        "interval": 30000000000,
        "conditions": ["[STATUS] == 200"]
      }
    ]
  }
}`
	code, body = runAdminV2Request(t, router, http.MethodPost, "/api/v1/admin/import", importDryRunPayload)
	if code != http.StatusOK {
		t.Fatalf("expected import dry run status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"endpointsCreated":1`) {
		t.Fatalf("expected import preview to indicate created endpoint, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/endpoints", "")
	if code != http.StatusOK {
		t.Fatalf("expected endpoints list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if strings.Contains(body, `"core_backend"`) {
		t.Fatalf("expected dry run import not to persist endpoint, got: %s", body)
	}

	importApplyPayload := strings.Replace(importDryRunPayload, `"dryRun": true`, `"dryRun": false`, 1)
	code, body = runAdminV2Request(t, router, http.MethodPost, "/api/v1/admin/import", importApplyPayload)
	if code != http.StatusOK {
		t.Fatalf("expected import apply status %d, got %d (%s)", http.StatusOK, code, body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/endpoints", "")
	if code != http.StatusOK {
		t.Fatalf("expected endpoints list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"core_backend"`) {
		t.Fatalf("expected applied import to persist endpoint, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/audit-logs?page=1&pageSize=10", "")
	if code != http.StatusOK {
		t.Fatalf("expected audit logs status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"items"`) {
		t.Fatalf("expected audit logs response to include items, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/audit-logs?page=1&pageSize=10&result=all", "")
	if code != http.StatusOK {
		t.Fatalf("expected audit logs (result=all) status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if strings.Contains(body, `"items":[]`) {
		t.Fatalf("expected audit logs with result=all not to be empty after operations, got: %s", body)
	}
}

func runAdminV2Request(t *testing.T, router *fiber.App, method, path, body string) (int, string) {
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

func loadAdminV2TestConfig(t *testing.T) *config.Config {
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
