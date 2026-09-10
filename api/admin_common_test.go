package api

import (
	"encoding/json"
	"testing"
)

func TestMarshalAuditSnapshotRedactsSensitiveValues(t *testing.T) {
	snapshot := map[string]any{
		"token": "super-secret-token",
		"headers": map[string]any{
			"Authorization": "Bearer super-secret",
			"X-Test":        "ok",
		},
		"nested": map[string]any{
			"clientSecret": "super-secret-client",
			"name":         "service-a",
		},
		"list": []any{
			map[string]any{"password": "p@ssw0rd"},
		},
		"webhook-url": "https://hooks.example.org/services/T00/B00/XXX",
		"webhookUrl":  "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=live-secret",
		"slack-webhook": "https://hooks.slack.com/services/T00/B00/XXX",
		"endpoint-url":  "https://status.example.org/health",
	}

	serialized := marshalAuditSnapshot(snapshot)
	if serialized == "" {
		t.Fatal("expected non-empty serialized audit snapshot")
	}

	decoded := map[string]any{}
	if err := json.Unmarshal([]byte(serialized), &decoded); err != nil {
		t.Fatalf("expected valid json snapshot, got error: %v", err)
	}

	if decoded["token"] != auditRedactedValue {
		t.Fatalf("expected token to be redacted, got %v", decoded["token"])
	}

	headers, ok := decoded["headers"].(map[string]any)
	if !ok {
		t.Fatalf("expected headers to be a map, got %T", decoded["headers"])
	}
	if headers["Authorization"] != auditRedactedValue {
		t.Fatalf("expected Authorization header to be redacted, got %v", headers["Authorization"])
	}
	if headers["X-Test"] != "ok" {
		t.Fatalf("expected non-sensitive header to remain unchanged, got %v", headers["X-Test"])
	}

	nested, ok := decoded["nested"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested to be a map, got %T", decoded["nested"])
	}
	if nested["clientSecret"] != auditRedactedValue {
		t.Fatalf("expected clientSecret to be redacted, got %v", nested["clientSecret"])
	}
	if nested["name"] != "service-a" {
		t.Fatalf("expected non-sensitive nested field to remain unchanged, got %v", nested["name"])
	}

	list, ok := decoded["list"].([]any)
	if !ok || len(list) != 1 {
		t.Fatalf("expected list to contain one element, got %T (%d)", decoded["list"], len(list))
	}
	listItem, ok := list[0].(map[string]any)
	if !ok {
		t.Fatalf("expected list item to be a map, got %T", list[0])
	}
	if listItem["password"] != auditRedactedValue {
		t.Fatalf("expected password in list item to be redacted, got %v", listItem["password"])
	}

	for _, key := range []string{"webhook-url", "webhookUrl", "slack-webhook"} {
		if decoded[key] != auditRedactedValue {
			t.Fatalf("expected %s to be redacted, got %v", key, decoded[key])
		}
	}
	if decoded["endpoint-url"] != "https://status.example.org/health" {
		t.Fatalf("expected bare url key to remain unchanged, got %v", decoded["endpoint-url"])
	}
}

func TestMarshalAuditSnapshotRedactsWeComWebhookPayload(t *testing.T) {
	// Mirrors notification create/update audit snapshots for WeCom/Slack.
	snapshot := map[string]any{
		"type": "custom",
		"name": "wecom-prod",
		"provider-override": map[string]any{
			"webhook-url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=corp-secret-key",
			"method":      "POST",
			"body":        `{"msgtype":"text","text":{"content":"alert"}}`,
		},
		"default-config": map[string]any{
			"webhookUrl": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=another-secret",
		},
	}

	serialized := marshalAuditSnapshot(snapshot)
	if serialized == "" {
		t.Fatal("expected non-empty serialized audit snapshot")
	}

	decoded := map[string]any{}
	if err := json.Unmarshal([]byte(serialized), &decoded); err != nil {
		t.Fatalf("expected valid json snapshot, got error: %v", err)
	}

	if decoded["name"] != "wecom-prod" {
		t.Fatalf("expected non-sensitive name to remain, got %v", decoded["name"])
	}

	override, ok := decoded["provider-override"].(map[string]any)
	if !ok {
		t.Fatalf("expected provider-override to be a map, got %T", decoded["provider-override"])
	}
	if override["webhook-url"] != auditRedactedValue {
		t.Fatalf("expected WeCom webhook-url to be redacted, got %v", override["webhook-url"])
	}
	if override["method"] != "POST" {
		t.Fatalf("expected method to remain unchanged, got %v", override["method"])
	}

	defaultConfig, ok := decoded["default-config"].(map[string]any)
	if !ok {
		t.Fatalf("expected default-config to be a map, got %T", decoded["default-config"])
	}
	if defaultConfig["webhookUrl"] != auditRedactedValue {
		t.Fatalf("expected webhookUrl to be redacted, got %v", defaultConfig["webhookUrl"])
	}
}
