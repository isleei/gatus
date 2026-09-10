package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAdminImportYAMLDryRunAndApply(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	router := New(cfg).Router()

	yamlBody := `endpoints:
  - name: yaml-imported
    group: core
    url: https://example.org/yaml
    interval: 30s
    conditions:
      - "[STATUS] == 200"
`

	// Dry-run via Content-Type: application/yaml + query overrides
	code, body := runAdminRequestWithContentType(t, router, http.MethodPost,
		"/api/v1/admin/import?entityType=endpoint&mode=merge&dryRun=true",
		yamlBody, "application/yaml")
	if code != http.StatusOK {
		t.Fatalf("expected YAML import dry-run status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"dryRun":true`) {
		t.Fatalf("expected dryRun=true in response, got: %s", body)
	}
	if !strings.Contains(body, `"endpointsCreated":1`) {
		t.Fatalf("expected endpointsCreated=1 in dry-run preview, got: %s", body)
	}
	if !strings.Contains(body, "Dry run completed") {
		t.Fatalf("expected dry-run message, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/endpoints", "")
	if code != http.StatusOK {
		t.Fatalf("expected endpoints list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if strings.Contains(body, `"core_yaml-imported"`) {
		t.Fatalf("expected dry-run YAML import not to persist endpoint, got: %s", body)
	}

	// Apply via Content-Type: application/x-yaml (also accepted — contains "yaml")
	code, body = runAdminRequestWithContentType(t, router, http.MethodPost,
		"/api/v1/admin/import?entityType=endpoint&mode=merge&dryRun=false",
		yamlBody, "application/x-yaml")
	if code != http.StatusOK {
		t.Fatalf("expected YAML import apply status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"dryRun":false`) {
		t.Fatalf("expected dryRun=false in apply response, got: %s", body)
	}
	if !strings.Contains(body, "Import applied") {
		t.Fatalf("expected apply message, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodGet, "/api/v1/admin/endpoints", "")
	if code != http.StatusOK {
		t.Fatalf("expected endpoints list status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"core_yaml-imported"`) {
		t.Fatalf("expected applied YAML import to persist endpoint, got: %s", body)
	}
}

func TestAdminImportYAMLWrappedRequestAndFormatQuery(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	router := New(cfg).Router()

	// Shape A: full ManagedImportRequest as YAML; format=yaml query without yaml Content-Type
	wrapped := `entityType: endpoint
mode: merge
dryRun: true
data:
  endpoints:
    - name: wrapped-yaml
      group: ops
      url: https://example.org/wrapped
      interval: 1m
      conditions:
        - "[STATUS] == 200"
`
	code, body := runAdminRequestWithContentType(t, router, http.MethodPost,
		"/api/v1/admin/import?format=yaml",
		wrapped, "text/plain")
	if code != http.StatusOK {
		t.Fatalf("expected wrapped YAML dry-run status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"endpointsCreated":1`) {
		t.Fatalf("expected endpointsCreated=1, got: %s", body)
	}
	if !strings.Contains(body, `"dryRun":true`) {
		t.Fatalf("expected dryRun=true, got: %s", body)
	}
}

func runAdminRequestWithContentType(t *testing.T, router *fiber.App, method, path, body, contentType string) (int, string) {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if len(body) > 0 && len(contentType) > 0 {
		request.Header.Set("Content-Type", contentType)
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
