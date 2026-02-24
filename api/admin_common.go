package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common"
	"github.com/TwiN/logr"
	"github.com/gofiber/fiber/v2"
)

const (
	adminListCacheTTL  = 5 * time.Second
	auditRedactedValue = "[REDACTED]"
)

var (
	adminAuditSensitiveKeySubstrings = []string{
		"accesskey",
		"apikey",
		"authorization",
		"bearer",
		"clientsecret",
		"credential",
		"passwd",
		"password",
		"privatekey",
		"secret",
		"session",
		"token",
	}
	adminAuditKeyNormalizer = strings.NewReplacer("_", "", "-", "", " ", "")
)

func clearAdminDerivedCache() {
	keys := cache.GetKeysByPattern("admin-*", 0)
	if len(keys) == 0 {
		return
	}
	_ = cache.DeleteAll(keys)
}

func buildAdminCacheKey(prefix string, c *fiber.Ctx) string {
	if c == nil {
		return "admin-" + prefix
	}
	return "admin-" + prefix + "-" + c.OriginalURL()
}

func writeAdminAudit(c *fiber.Ctx, cfg *config.Config, action, entityType, entityKey string, before, after interface{}, operationErr error) {
	entry := &common.AdminAuditLogEntry{
		Actor:      "anonymous",
		Action:     strings.TrimSpace(action),
		EntityType: strings.TrimSpace(entityType),
		EntityKey:  strings.TrimSpace(entityKey),
		Result:     "success",
		RequestID:  "",
		Timestamp:  time.Now().UTC(),
	}
	if c != nil {
		entry.RequestID = strings.TrimSpace(c.Get("X-Request-ID"))
		if len(entry.RequestID) == 0 {
			entry.RequestID = strings.TrimSpace(c.Get("X-Request-Id"))
		}
	}
	if cfg != nil && cfg.Security != nil {
		entry.Actor = cfg.Security.GetActor(c)
	}
	if operationErr != nil {
		entry.Result = "failure"
		entry.Error = operationErr.Error()
	}
	entry.Before = marshalAuditSnapshot(before)
	entry.After = marshalAuditSnapshot(after)
	if err := store.Get().InsertAdminAuditLog(entry); err != nil {
		logr.Errorf("[api.writeAdminAudit] Failed to persist audit log action=%s entityType=%s entityKey=%s: %s", action, entityType, entityKey, err.Error())
	}
}

func marshalAuditSnapshot(value interface{}) string {
	if value == nil {
		return ""
	}
	serialized, err := json.Marshal(redactAuditSnapshot(value))
	if err != nil {
		return ""
	}
	return string(serialized)
}

func redactAuditSnapshot(value interface{}) interface{} {
	serialized, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var decoded interface{}
	if err := json.Unmarshal(serialized, &decoded); err != nil {
		return value
	}
	return redactAuditDecodedValue(decoded, "")
}

func redactAuditDecodedValue(value interface{}, key string) interface{} {
	if isAdminAuditSensitiveKey(key) {
		return auditRedactedValue
	}
	switch typed := value.(type) {
	case map[string]interface{}:
		redacted := make(map[string]interface{}, len(typed))
		for mapKey, mapValue := range typed {
			if isAdminAuditSensitiveKey(mapKey) {
				redacted[mapKey] = auditRedactedValue
				continue
			}
			redacted[mapKey] = redactAuditDecodedValue(mapValue, mapKey)
		}
		return redacted
	case []interface{}:
		redacted := make([]interface{}, len(typed))
		for i, entry := range typed {
			redacted[i] = redactAuditDecodedValue(entry, key)
		}
		return redacted
	default:
		return value
	}
}

func isAdminAuditSensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.TrimSpace(key))
	if normalized == "" {
		return false
	}
	normalized = adminAuditKeyNormalizer.Replace(normalized)
	for _, candidate := range adminAuditSensitiveKeySubstrings {
		if strings.Contains(normalized, candidate) {
			return true
		}
	}
	return false
}
