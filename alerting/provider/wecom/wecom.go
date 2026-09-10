package wecom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/client"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"gopkg.in/yaml.v3"
)

var (
	ErrWebhookURLNotSet       = errors.New("webhook-url not set")
	ErrDuplicateGroupOverride = errors.New("duplicate group override")
)

type Config struct {
	WebhookURL string `yaml:"webhook-url"`
	Title      string `yaml:"title,omitempty"`

	// TextTriggered is an optional markdown body template used when an alert is triggered.
	// Supported placeholders: [ENDPOINT], [ENDPOINT_NAME], [ENDPOINT_GROUP], [ALERT_DESCRIPTION],
	// [FAILURE_COUNT], [SUCCESS_COUNT], [RESULT_CONDITIONS], [RESULT_ERRORS].
	// When empty, the historical English default body is used unchanged.
	TextTriggered string `yaml:"text-triggered,omitempty"`

	// TextResolved is an optional markdown body template used when an alert is resolved.
	// Same placeholders as TextTriggered. When empty, the historical English default body is used unchanged.
	TextResolved string `yaml:"text-resolved,omitempty"`
}

func (cfg *Config) Validate() error {
	if len(cfg.WebhookURL) == 0 {
		return ErrWebhookURLNotSet
	}
	return nil
}

func (cfg *Config) Merge(override *Config) {
	if len(override.WebhookURL) > 0 {
		cfg.WebhookURL = override.WebhookURL
	}
	if len(override.Title) > 0 {
		cfg.Title = override.Title
	}
	if len(override.TextTriggered) > 0 {
		cfg.TextTriggered = override.TextTriggered
	}
	if len(override.TextResolved) > 0 {
		cfg.TextResolved = override.TextResolved
	}
}

// AlertProvider is the configuration necessary for sending an alert using WeCom's webhook robot.
type AlertProvider struct {
	DefaultConfig Config `yaml:",inline"`

	// DefaultAlert is the default alert configuration to use for endpoints with an alert of the appropriate type
	DefaultAlert *alert.Alert `yaml:"default-alert,omitempty"`

	// Overrides is a list of Override that may be prioritized over the default configuration
	Overrides []Override `yaml:"overrides,omitempty"`
}

// Override is a case under which the default integration is overridden
type Override struct {
	Group  string `yaml:"group"`
	Config `yaml:",inline"`
}

// Validate the provider's configuration
func (provider *AlertProvider) Validate() error {
	registeredGroups := make(map[string]bool)
	if provider.Overrides != nil {
		for _, override := range provider.Overrides {
			if isAlreadyRegistered := registeredGroups[override.Group]; isAlreadyRegistered || override.Group == "" {
				return ErrDuplicateGroupOverride
			}
			registeredGroups[override.Group] = true
		}
	}
	return provider.DefaultConfig.Validate()
}

// Send an alert using the provider
func (provider *AlertProvider) Send(ep *endpoint.Endpoint, alert *alert.Alert, result *endpoint.Result, resolved bool) error {
	cfg, err := provider.GetConfig(ep.Group, alert)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(provider.buildRequestBody(cfg, ep, alert, result, resolved))
	request, err := http.NewRequest(http.MethodPost, cfg.WebhookURL, buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.GetHTTPClient(nil).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("call to provider alert returned status code %d: %s", response.StatusCode, string(body))
	}
	return nil
}

type Body struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

// buildRequestBody builds the request body for the provider
func (provider *AlertProvider) buildRequestBody(cfg *Config, ep *endpoint.Endpoint, alert *alert.Alert, result *endpoint.Result, resolved bool) []byte {
	title := cfg.Title
	if len(title) == 0 {
		title = "Gatus"
	}
	var bodyContent string
	if resolved {
		if len(cfg.TextResolved) > 0 {
			bodyContent = provider.renderTextTemplate(cfg.TextResolved, ep, alert, result)
		} else {
			bodyContent = provider.defaultBodyContent(ep, alert, result, true)
		}
	} else {
		if len(cfg.TextTriggered) > 0 {
			bodyContent = provider.renderTextTemplate(cfg.TextTriggered, ep, alert, result)
		} else {
			bodyContent = provider.defaultBodyContent(ep, alert, result, false)
		}
	}
	body := Body{
		MsgType: "markdown",
		Markdown: Markdown{
			Content: fmt.Sprintf("**%s**\n%s", title, bodyContent),
		},
	}
	bodyAsJSON, _ := json.Marshal(body)
	return bodyAsJSON
}

// defaultBodyContent preserves the historical English WeCom body (message + description + conditions).
func (provider *AlertProvider) defaultBodyContent(ep *endpoint.Endpoint, alert *alert.Alert, result *endpoint.Result, resolved bool) string {
	var message string
	if resolved {
		message = fmt.Sprintf("An alert for **%s** has been resolved after passing successfully %d time(s) in a row", ep.DisplayName(), alert.SuccessThreshold)
	} else {
		message = fmt.Sprintf("An alert for **%s** has been triggered due to having failed %d time(s) in a row", ep.DisplayName(), alert.FailureThreshold)
	}
	description := ""
	if alertDescription := alert.GetDescription(); len(alertDescription) > 0 {
		description = "\nDescription:\n> " + alertDescription
	}
	conditionResults := ""
	if len(result.ConditionResults) > 0 {
		conditionResults = "\nCondition results:\n"
		for _, conditionResult := range result.ConditionResults {
			prefix := "[FAIL]"
			if conditionResult.Success {
				prefix = "[PASS]"
			}
			conditionResults += fmt.Sprintf("%s %s\n", prefix, conditionResult.Condition)
		}
	}
	return message + description + conditionResults
}

// renderTextTemplate replaces supported placeholders in a custom text template.
// Values are inserted as-is (no extra escaping) to match other providers; callers should avoid
// embedding untrusted markdown in endpoint names/descriptions if that is a concern.
func (provider *AlertProvider) renderTextTemplate(tmpl string, ep *endpoint.Endpoint, alert *alert.Alert, result *endpoint.Result) string {
	message := tmpl
	message = strings.ReplaceAll(message, "[ENDPOINT]", ep.DisplayName())
	message = strings.ReplaceAll(message, "[ENDPOINT_NAME]", ep.Name)
	message = strings.ReplaceAll(message, "[ENDPOINT_GROUP]", ep.Group)
	message = strings.ReplaceAll(message, "[ALERT_DESCRIPTION]", alert.GetDescription())
	message = strings.ReplaceAll(message, "[FAILURE_COUNT]", strconv.Itoa(alert.FailureThreshold))
	message = strings.ReplaceAll(message, "[SUCCESS_COUNT]", strconv.Itoa(alert.SuccessThreshold))
	if strings.Contains(message, "[RESULT_CONDITIONS]") {
		message = strings.ReplaceAll(message, "[RESULT_CONDITIONS]", formatConditionResults(result))
	}
	if strings.Contains(message, "[RESULT_ERRORS]") {
		message = strings.ReplaceAll(message, "[RESULT_ERRORS]", strings.Join(result.Errors, ", "))
	}
	return message
}

func formatConditionResults(result *endpoint.Result) string {
	if result == nil || len(result.ConditionResults) == 0 {
		return ""
	}
	var b strings.Builder
	for i, conditionResult := range result.ConditionResults {
		prefix := "[FAIL]"
		if conditionResult.Success {
			prefix = "[PASS]"
		}
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(prefix)
		b.WriteByte(' ')
		b.WriteString(conditionResult.Condition)
	}
	return b.String()
}

// GetDefaultAlert returns the provider's default alert configuration
func (provider *AlertProvider) GetDefaultAlert() *alert.Alert {
	return provider.DefaultAlert
}

// GetConfig returns the configuration for the provider with the overrides applied
func (provider *AlertProvider) GetConfig(group string, alert *alert.Alert) (*Config, error) {
	cfg := provider.DefaultConfig
	// Handle group overrides
	if provider.Overrides != nil {
		for _, override := range provider.Overrides {
			if group == override.Group {
				cfg.Merge(&override.Config)
				break
			}
		}
	}
	// Handle alert overrides
	if len(alert.ProviderOverride) != 0 {
		overrideConfig := Config{}
		if err := yaml.Unmarshal(alert.ProviderOverrideAsBytes(), &overrideConfig); err != nil {
			return nil, err
		}
		cfg.Merge(&overrideConfig)
	}
	// Validate the configuration
	err := cfg.Validate()
	return &cfg, err
}

// ValidateOverrides validates the alert's provider override and, if present, the group override
func (provider *AlertProvider) ValidateOverrides(group string, alert *alert.Alert) error {
	_, err := provider.GetConfig(group, alert)
	return err
}
