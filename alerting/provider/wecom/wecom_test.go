package wecom

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/client"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/test"
)

func TestAlertProvider_Validate(t *testing.T) {
	invalidProvider := AlertProvider{DefaultConfig: Config{WebhookURL: ""}}
	if err := invalidProvider.Validate(); err == nil {
		t.Error("provider shouldn't have been valid")
	}
	validProvider := AlertProvider{DefaultConfig: Config{WebhookURL: "https://example.com/webhook"}}
	if err := validProvider.Validate(); err != nil {
		t.Error("provider should've been valid")
	}
}

func TestAlertProvider_ValidateWithOverride(t *testing.T) {
	providerWithInvalidOverrideGroup := AlertProvider{
		DefaultConfig: Config{WebhookURL: "https://example.com/webhook"},
		Overrides: []Override{
			{
				Config: Config{WebhookURL: "https://example.com/override"},
				Group:  "",
			},
		},
	}
	if err := providerWithInvalidOverrideGroup.Validate(); err == nil {
		t.Error("provider Group shouldn't have been valid")
	}
	providerWithDuplicateOverrideGroup := AlertProvider{
		DefaultConfig: Config{WebhookURL: "https://example.com/webhook"},
		Overrides: []Override{
			{
				Config: Config{WebhookURL: "https://example.com/override-1"},
				Group:  "group",
			},
			{
				Config: Config{WebhookURL: "https://example.com/override-2"},
				Group:  "group",
			},
		},
	}
	if err := providerWithDuplicateOverrideGroup.Validate(); err == nil {
		t.Error("provider duplicate Group shouldn't have been valid")
	}
	providerWithValidOverride := AlertProvider{
		DefaultConfig: Config{WebhookURL: "https://example.com/webhook"},
		Overrides: []Override{
			{
				Config: Config{WebhookURL: "https://example.com/group-webhook"},
				Group:  "group",
			},
		},
	}
	if err := providerWithValidOverride.Validate(); err != nil {
		t.Error("provider should've been valid")
	}
}

func TestAlertProvider_Send(t *testing.T) {
	defer client.InjectHTTPClient(nil)
	description := "description-1"
	scenarios := []struct {
		Name             string
		Resolved         bool
		MockRoundTripper test.MockRoundTripper
		ExpectedError    bool
	}{
		{
			Name:     "triggered",
			Resolved: false,
			MockRoundTripper: test.MockRoundTripper(func(r *http.Request) *http.Response {
				return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
			}),
			ExpectedError: false,
		},
		{
			Name:     "triggered-error",
			Resolved: false,
			MockRoundTripper: test.MockRoundTripper(func(r *http.Request) *http.Response {
				return &http.Response{StatusCode: http.StatusInternalServerError, Body: http.NoBody}
			}),
			ExpectedError: true,
		},
		{
			Name:     "resolved",
			Resolved: true,
			MockRoundTripper: test.MockRoundTripper(func(r *http.Request) *http.Response {
				return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}
			}),
			ExpectedError: false,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			client.InjectHTTPClient(&http.Client{Transport: scenario.MockRoundTripper})
			err := (&AlertProvider{
				DefaultConfig: Config{WebhookURL: "https://example.com/webhook"},
			}).Send(
				&endpoint.Endpoint{Name: "endpoint-name"},
				&alert.Alert{
					Description:      &description,
					SuccessThreshold: 5,
					FailureThreshold: 3,
				},
				&endpoint.Result{
					ConditionResults: []*endpoint.ConditionResult{
						{Condition: "[CONNECTED] == true", Success: scenario.Resolved},
						{Condition: "[STATUS] == 200", Success: scenario.Resolved},
					},
				},
				scenario.Resolved,
			)
			if scenario.ExpectedError && err == nil {
				t.Error("expected error, got none")
			}
			if !scenario.ExpectedError && err != nil {
				t.Error("expected no error, got", err.Error())
			}
		})
	}
}

func TestAlertProvider_buildRequestBody(t *testing.T) {
	description := "description-1"
	scenarios := []struct {
		Name             string
		Title            string
		Resolved         bool
		NoConditions     bool
		ExpectedContains []string
	}{
		{
			Name:         "triggered",
			Title:        "",
			Resolved:     false,
			NoConditions: false,
			ExpectedContains: []string{
				"**Gatus**",
				"has been triggered due to having failed 3 time(s) in a row",
				"Description:",
				"Condition results:",
			},
		},
		{
			Name:         "resolved",
			Title:        "",
			Resolved:     true,
			NoConditions: false,
			ExpectedContains: []string{
				"**Gatus**",
				"has been resolved after passing successfully 5 time(s) in a row",
				"Condition results:",
			},
		},
		{
			Name:         "resolved-with-custom-title-and-no-conditions",
			Title:        "custom-title",
			Resolved:     true,
			NoConditions: true,
			ExpectedContains: []string{
				"**custom-title**",
				"has been resolved after passing successfully 5 time(s) in a row",
			},
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			var conditionResults []*endpoint.ConditionResult
			if !scenario.NoConditions {
				conditionResults = []*endpoint.ConditionResult{
					{Condition: "[CONNECTED] == true", Success: scenario.Resolved},
					{Condition: "[STATUS] == 200", Success: scenario.Resolved},
				}
			}
			bodyBytes := (&AlertProvider{}).buildRequestBody(
				&Config{WebhookURL: "https://example.com/webhook", Title: scenario.Title},
				&endpoint.Endpoint{Name: "endpoint-name"},
				&alert.Alert{
					Description:      &description,
					SuccessThreshold: 5,
					FailureThreshold: 3,
				},
				&endpoint.Result{ConditionResults: conditionResults},
				scenario.Resolved,
			)
			var body Body
			if err := json.Unmarshal(bodyBytes, &body); err != nil {
				t.Fatal("expected body to be valid JSON, got error:", err.Error())
			}
			if body.MsgType != "markdown" {
				t.Fatalf("expected msgtype markdown, got %s", body.MsgType)
			}
			for _, expected := range scenario.ExpectedContains {
				if !strings.Contains(body.Markdown.Content, expected) {
					t.Errorf("expected markdown content to contain %q, got %q", expected, body.Markdown.Content)
				}
			}
		})
	}
}

func TestAlertProvider_buildRequestBody_defaultExact(t *testing.T) {
	description := "description-1"
	ep := &endpoint.Endpoint{Name: "endpoint-name", Group: "core"}
	alertCfg := &alert.Alert{
		Description:      &description,
		SuccessThreshold: 5,
		FailureThreshold: 3,
	}
	result := &endpoint.Result{
		ConditionResults: []*endpoint.ConditionResult{
			{Condition: "[CONNECTED] == true", Success: false},
			{Condition: "[STATUS] == 200", Success: false},
		},
	}
	bodyBytes := (&AlertProvider{}).buildRequestBody(
		&Config{WebhookURL: "https://example.com/webhook"},
		ep,
		alertCfg,
		result,
		false,
	)
	var body Body
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatal(err)
	}
	want := "**Gatus**\nAn alert for **core/endpoint-name** has been triggered due to having failed 3 time(s) in a row\nDescription:\n> description-1\nCondition results:\n[FAIL] [CONNECTED] == true\n[FAIL] [STATUS] == 200\n"
	if body.Markdown.Content != want {
		t.Errorf("default triggered body mismatch\nwant: %q\ngot:  %q", want, body.Markdown.Content)
	}

	resolvedBytes := (&AlertProvider{}).buildRequestBody(
		&Config{WebhookURL: "https://example.com/webhook"},
		ep,
		alertCfg,
		&endpoint.Result{
			ConditionResults: []*endpoint.ConditionResult{
				{Condition: "[CONNECTED] == true", Success: true},
				{Condition: "[STATUS] == 200", Success: true},
			},
		},
		true,
	)
	var resolvedBody Body
	if err := json.Unmarshal(resolvedBytes, &resolvedBody); err != nil {
		t.Fatal(err)
	}
	wantResolved := "**Gatus**\nAn alert for **core/endpoint-name** has been resolved after passing successfully 5 time(s) in a row\nDescription:\n> description-1\nCondition results:\n[PASS] [CONNECTED] == true\n[PASS] [STATUS] == 200\n"
	if resolvedBody.Markdown.Content != wantResolved {
		t.Errorf("default resolved body mismatch\nwant: %q\ngot:  %q", wantResolved, resolvedBody.Markdown.Content)
	}
}

func TestAlertProvider_buildRequestBody_customTemplates(t *testing.T) {
	description := "超时"
	ep := &endpoint.Endpoint{Name: "checkout", Group: "critical"}
	alertCfg := &alert.Alert{
		Description:      &description,
		SuccessThreshold: 2,
		FailureThreshold: 3,
	}
	result := &endpoint.Result{
		ConditionResults: []*endpoint.ConditionResult{
			{Condition: "[STATUS] == 200", Success: false},
		},
		Errors: []string{"connection refused"},
	}
	cfg := &Config{
		WebhookURL: "https://example.com/webhook",
		Title:      "监控告警",
		TextTriggered: "告警触发: **[ENDPOINT]**\n名称: [ENDPOINT_NAME]\n组: [ENDPOINT_GROUP]\n失败阈值: [FAILURE_COUNT]\n描述: [ALERT_DESCRIPTION]\n条件:\n[RESULT_CONDITIONS]\n错误: [RESULT_ERRORS]",
		TextResolved:  "告警恢复: **[ENDPOINT]** (连续成功 [SUCCESS_COUNT] 次)\n描述: [ALERT_DESCRIPTION]",
	}
	triggered := (&AlertProvider{}).buildRequestBody(cfg, ep, alertCfg, result, false)
	var triggeredBody Body
	if err := json.Unmarshal(triggered, &triggeredBody); err != nil {
		t.Fatal(err)
	}
	wantTriggered := "**监控告警**\n告警触发: **critical/checkout**\n名称: checkout\n组: critical\n失败阈值: 3\n描述: 超时\n条件:\n[FAIL] [STATUS] == 200\n错误: connection refused"
	if triggeredBody.Markdown.Content != wantTriggered {
		t.Errorf("custom triggered mismatch\nwant: %q\ngot:  %q", wantTriggered, triggeredBody.Markdown.Content)
	}

	resolved := (&AlertProvider{}).buildRequestBody(cfg, ep, alertCfg, result, true)
	var resolvedBody Body
	if err := json.Unmarshal(resolved, &resolvedBody); err != nil {
		t.Fatal(err)
	}
	wantResolved := "**监控告警**\n告警恢复: **critical/checkout** (连续成功 2 次)\n描述: 超时"
	if resolvedBody.Markdown.Content != wantResolved {
		t.Errorf("custom resolved mismatch\nwant: %q\ngot:  %q", wantResolved, resolvedBody.Markdown.Content)
	}

	// Custom triggered should not fall back to English default fragments
	if strings.Contains(triggeredBody.Markdown.Content, "has been triggered") {
		t.Error("custom template should not include English default message")
	}
}

func TestConfig_Merge_textTemplates(t *testing.T) {
	cfg := Config{
		WebhookURL:    "https://example.com/default",
		Title:         "default",
		TextTriggered: "triggered-default",
		TextResolved:  "resolved-default",
	}
	cfg.Merge(&Config{
		TextTriggered: "triggered-override",
		TextResolved:  "resolved-override",
	})
	if cfg.TextTriggered != "triggered-override" {
		t.Errorf("expected triggered override, got %q", cfg.TextTriggered)
	}
	if cfg.TextResolved != "resolved-override" {
		t.Errorf("expected resolved override, got %q", cfg.TextResolved)
	}
	if cfg.WebhookURL != "https://example.com/default" {
		t.Errorf("webhook-url should remain default, got %q", cfg.WebhookURL)
	}
}

func TestAlertProvider_GetDefaultAlert(t *testing.T) {
	if (&AlertProvider{DefaultAlert: &alert.Alert{}}).GetDefaultAlert() == nil {
		t.Error("expected default alert to be not nil")
	}
	if (&AlertProvider{DefaultAlert: nil}).GetDefaultAlert() != nil {
		t.Error("expected default alert to be nil")
	}
}

func TestAlertProvider_GetConfig(t *testing.T) {
	scenarios := []struct {
		Name           string
		Provider       AlertProvider
		InputGroup     string
		InputAlert     alert.Alert
		ExpectedOutput Config
		ExpectedError  bool
	}{
		{
			Name: "provider-no-override",
			Provider: AlertProvider{
				DefaultConfig: Config{WebhookURL: "https://example.com/default", Title: "default"},
			},
			InputGroup:     "",
			InputAlert:     alert.Alert{},
			ExpectedOutput: Config{WebhookURL: "https://example.com/default", Title: "default"},
		},
		{
			Name: "provider-with-group-override",
			Provider: AlertProvider{
				DefaultConfig: Config{WebhookURL: "https://example.com/default", Title: "default"},
				Overrides: []Override{
					{
						Group:  "group",
						Config: Config{WebhookURL: "https://example.com/group", Title: "group-title"},
					},
				},
			},
			InputGroup:     "group",
			InputAlert:     alert.Alert{},
			ExpectedOutput: Config{WebhookURL: "https://example.com/group", Title: "group-title"},
		},
		{
			Name: "provider-with-alert-override",
			Provider: AlertProvider{
				DefaultConfig: Config{WebhookURL: "https://example.com/default", Title: "default"},
			},
			InputGroup:     "",
			InputAlert:     alert.Alert{ProviderOverride: map[string]any{"webhook-url": "https://example.com/alert", "title": "alert-title"}},
			ExpectedOutput: Config{WebhookURL: "https://example.com/alert", Title: "alert-title"},
		},
		{
			Name: "provider-with-empty-alert-override-should-not-clear-default",
			Provider: AlertProvider{
				DefaultConfig: Config{WebhookURL: "https://example.com/default", Title: "default"},
			},
			InputGroup:     "",
			InputAlert:     alert.Alert{ProviderOverride: map[string]any{"webhook-url": ""}},
			ExpectedOutput: Config{WebhookURL: "https://example.com/default", Title: "default"},
		},
		{
			Name: "provider-with-text-template-alert-override",
			Provider: AlertProvider{
				DefaultConfig: Config{
					WebhookURL:    "https://example.com/default",
					TextTriggered: "default-triggered",
					TextResolved:  "default-resolved",
				},
			},
			InputGroup: "",
			InputAlert: alert.Alert{ProviderOverride: map[string]any{
				"text-triggered": "alert-triggered",
				"text-resolved":  "alert-resolved",
			}},
			ExpectedOutput: Config{
				WebhookURL:    "https://example.com/default",
				TextTriggered: "alert-triggered",
				TextResolved:  "alert-resolved",
			},
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			cfg, err := scenario.Provider.GetConfig(scenario.InputGroup, &scenario.InputAlert)
			if scenario.ExpectedError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal("expected no error, got", err.Error())
			}
			if cfg.WebhookURL != scenario.ExpectedOutput.WebhookURL {
				t.Errorf("expected webhook-url %s, got %s", scenario.ExpectedOutput.WebhookURL, cfg.WebhookURL)
			}
			if cfg.Title != scenario.ExpectedOutput.Title {
				t.Errorf("expected title %s, got %s", scenario.ExpectedOutput.Title, cfg.Title)
			}
			if cfg.TextTriggered != scenario.ExpectedOutput.TextTriggered {
				t.Errorf("expected text-triggered %s, got %s", scenario.ExpectedOutput.TextTriggered, cfg.TextTriggered)
			}
			if cfg.TextResolved != scenario.ExpectedOutput.TextResolved {
				t.Errorf("expected text-resolved %s, got %s", scenario.ExpectedOutput.TextResolved, cfg.TextResolved)
			}
		})
	}
}
