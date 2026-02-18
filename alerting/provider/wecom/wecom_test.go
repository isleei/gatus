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
		})
	}
}
