package api

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/TwiN/gatus/v5/security"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common"
	"github.com/gofiber/fiber/v2"
)

func TestDeleteAdminAuditLogsHTTP_HappyPathDaysAndMaxAge(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	router := New(cfg).Router()

	old := time.Now().UTC().Add(-48 * time.Hour)
	recent := time.Now().UTC().Add(-1 * time.Hour)
	if err := store.Get().InsertAdminAuditLog(&common.AdminAuditLogEntry{
		Actor: "tester", Action: "old-cleanup-seed", EntityType: "endpoint", Result: "success", Timestamp: old,
	}); err != nil {
		t.Fatalf("unexpected insert error: %v", err)
	}
	if err := store.Get().InsertAdminAuditLog(&common.AdminAuditLogEntry{
		Actor: "tester", Action: "recent-cleanup-seed", EntityType: "endpoint", Result: "success", Timestamp: recent,
	}); err != nil {
		t.Fatalf("unexpected insert error: %v", err)
	}

	code, body := runAdminV2Request(t, router, http.MethodDelete, "/api/v1/admin/audit-logs?days=1", "")
	if code != http.StatusOK {
		t.Fatalf("expected DELETE days=1 status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"maxAge":"24h0m0s"`) {
		t.Fatalf("expected maxAge 24h for days=1, got: %s", body)
	}
	if strings.Contains(body, `"deleted":0`) {
		t.Fatalf("expected at least one audit log deleted for days=1, got: %s", body)
	}

	// Seed another old row for maxAge path
	if err := store.Get().InsertAdminAuditLog(&common.AdminAuditLogEntry{
		Actor: "tester", Action: "old-maxage-seed", EntityType: "endpoint", Result: "success",
		Timestamp: time.Now().UTC().Add(-3 * time.Hour),
	}); err != nil {
		t.Fatalf("unexpected insert error: %v", err)
	}
	code, body = runAdminV2Request(t, router, http.MethodDelete, "/api/v1/admin/audit-logs?maxAge=2h", "")
	if code != http.StatusOK {
		t.Fatalf("expected DELETE maxAge=2h status %d, got %d (%s)", http.StatusOK, code, body)
	}
	if !strings.Contains(body, `"maxAge":"2h0m0s"`) {
		t.Fatalf("expected maxAge 2h in response, got: %s", body)
	}
	if strings.Contains(body, `"deleted":0`) {
		t.Fatalf("expected at least one audit log deleted for maxAge=2h, got: %s", body)
	}
}

func TestDeleteAdminAuditLogsHTTP_BadRequestMissingParams(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	router := New(cfg).Router()

	code, body := runAdminV2Request(t, router, http.MethodDelete, "/api/v1/admin/audit-logs", "")
	if code != http.StatusBadRequest {
		t.Fatalf("expected status %d when maxAge/days missing, got %d (%s)", http.StatusBadRequest, code, body)
	}
	if !strings.Contains(body, "maxAge or days") {
		t.Fatalf("expected helpful error about maxAge/days, got: %s", body)
	}

	code, body = runAdminV2Request(t, router, http.MethodDelete, "/api/v1/admin/audit-logs?maxAge=not-a-duration", "")
	if code != http.StatusBadRequest {
		t.Fatalf("expected status %d for invalid maxAge, got %d (%s)", http.StatusBadRequest, code, body)
	}
}

func TestDeleteAdminAuditLogsHTTP_RequiresAuth(t *testing.T) {
	cfg := loadAdminV2TestConfig(t)
	cfg.Security = &security.Config{
		Basic: &security.BasicConfig{
			Username:                        "john.doe",
			PasswordBcryptHashBase64Encoded: "JDJhJDA4JDFoRnpPY1hnaFl1OC9ISlFsa21VS09wOGlPU1ZOTDlHZG1qeTFvb3dIckRBUnlHUmNIRWlT", // hunter2
		},
	}
	router := New(cfg).Router()

	code, body := runAdminAuditDeleteRequest(t, router, "/api/v1/admin/audit-logs?days=30", "")
	if code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized without credentials, got %d (%s)", code, body)
	}

	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("john.doe:hunter2"))
	code, body = runAdminAuditDeleteRequest(t, router, "/api/v1/admin/audit-logs?days=30", auth)
	if code != http.StatusOK {
		t.Fatalf("expected OK with valid basic auth, got %d (%s)", code, body)
	}
	if !strings.Contains(body, `"deleted":`) {
		t.Fatalf("expected deleted field with auth, got: %s", body)
	}
}

func runAdminAuditDeleteRequest(t *testing.T, router *fiber.App, path, authorization string) (int, string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodDelete, path, http.NoBody)
	if len(authorization) > 0 {
		request.Header.Set("Authorization", authorization)
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
