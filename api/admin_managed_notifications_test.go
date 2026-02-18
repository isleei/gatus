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

func TestManagedNotificationsCRUD(t *testing.T) {
	cfg := loadManagedNotificationsTestConfig(t, false)
	router := New(cfg).Router()

	getCode, getBody := runManagedNotificationRequest(t, router, http.MethodGet, "/api/v1/admin/notifications", "")
	if getCode != http.StatusOK {
		t.Fatalf("expected GET code %d, got %d", http.StatusOK, getCode)
	}
	if !strings.Contains(getBody, `"type":"slack"`) {
		t.Fatalf("expected response to include slack notification type, got: %s", getBody)
	}

	slackPayload := `{
  "webhook-url": "https://example.org/slack-webhook",
  "title": "Infra Alerts"
}`
	putCode, putBody := runManagedNotificationRequest(t, router, http.MethodPut, "/api/v1/admin/notifications/slack", slackPayload)
	if putCode != http.StatusOK {
		t.Fatalf("expected PUT code %d, got %d (%s)", http.StatusOK, putCode, putBody)
	}
	if !strings.Contains(putBody, `"configured":true`) {
		t.Fatalf("expected saved notification to be configured, got: %s", putBody)
	}

	overlayAfterPut, err := config.ReadManagedOverlay(cfg.LoadedConfigPath())
	if err != nil {
		t.Fatalf("unexpected overlay read error: %v", err)
	}
	if overlayAfterPut == nil || overlayAfterPut.Alerting == nil || overlayAfterPut.Alerting.Slack == nil {
		t.Fatalf("expected overlay to include slack notification, got %+v", overlayAfterPut)
	}
	if overlayAfterPut.Alerting.Slack.DefaultConfig.WebhookURL != "https://example.org/slack-webhook" {
		t.Fatalf("unexpected webhook URL, got %s", overlayAfterPut.Alerting.Slack.DefaultConfig.WebhookURL)
	}

	deleteCode, deleteBody := runManagedNotificationRequest(t, router, http.MethodDelete, "/api/v1/admin/notifications/slack", "")
	if deleteCode != http.StatusNoContent {
		t.Fatalf("expected DELETE code %d, got %d (%s)", http.StatusNoContent, deleteCode, deleteBody)
	}

	overlayAfterDelete, err := config.ReadManagedOverlay(cfg.LoadedConfigPath())
	if err != nil {
		t.Fatalf("unexpected overlay read error: %v", err)
	}
	if overlayAfterDelete == nil {
		t.Fatal("expected overlay to still exist")
	}
	if overlayAfterDelete.Alerting != nil && overlayAfterDelete.Alerting.Slack != nil {
		t.Fatalf("expected slack notification to be removed, got %+v", overlayAfterDelete.Alerting.Slack)
	}
}

func TestManagedNotificationPutInvalidProviderConfig(t *testing.T) {
	cfg := loadManagedNotificationsTestConfig(t, false)
	router := New(cfg).Router()

	putCode, putBody := runManagedNotificationRequest(t, router, http.MethodPut, "/api/v1/admin/notifications/telegram", `{}`)
	if putCode != http.StatusBadRequest {
		t.Fatalf("expected PUT code %d, got %d (%s)", http.StatusBadRequest, putCode, putBody)
	}
	if !strings.Contains(putBody, "token not set") {
		t.Fatalf("expected token validation error, got: %s", putBody)
	}
}

func TestManagedNotificationDeleteConflictWhenInUse(t *testing.T) {
	cfg := loadManagedNotificationsTestConfig(t, true)
	router := New(cfg).Router()

	deleteCode, deleteBody := runManagedNotificationRequest(t, router, http.MethodDelete, "/api/v1/admin/notifications/slack", "")
	if deleteCode != http.StatusConflict {
		t.Fatalf("expected DELETE code %d, got %d (%s)", http.StatusConflict, deleteCode, deleteBody)
	}
	if !strings.Contains(deleteBody, "still referenced") {
		t.Fatalf("expected conflict reason about references, got: %s", deleteBody)
	}
}

func runManagedNotificationRequest(t *testing.T, router *fiber.App, method, path, body string) (int, string) {
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

func loadManagedNotificationsTestConfig(t *testing.T, withSlackAlertReference bool) *config.Config {
	t.Helper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	alertSection := ""
	if withSlackAlertReference {
		alertSection = `
    alerts:
      - type: slack
`
	}
	configYAML := `endpoints:
  - name: frontend
    group: core
    url: https://example.org/health
    interval: 30s
    conditions:
      - "[STATUS] == 200"` + alertSection + `
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
