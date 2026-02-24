package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/alerting/alert"
)

type ManagedAlertPayload struct {
	Type                    string         `json:"type"`
	Enabled                 *bool          `json:"enabled,omitempty"`
	FailureThreshold        int            `json:"failureThreshold,omitempty"`
	SuccessThreshold        int            `json:"successThreshold,omitempty"`
	MinimumReminderInterval string         `json:"minimumReminderInterval,omitempty"`
	Description             *string        `json:"description,omitempty"`
	SendOnResolved          *bool          `json:"sendOnResolved,omitempty"`
	ProviderOverride        map[string]any `json:"providerOverride,omitempty"`
}

func buildManagedAlertsFromPayload(payload []ManagedAlertPayload) ([]*alert.Alert, error) {
	alerts := make([]*alert.Alert, 0, len(payload))
	for _, entry := range payload {
		typeName := strings.TrimSpace(strings.ToLower(entry.Type))
		if len(typeName) == 0 {
			continue
		}
		minimumReminderInterval, err := parseManagedAlertDuration(entry.MinimumReminderInterval)
		if err != nil {
			return nil, err
		}
		parsed := &alert.Alert{
			Type:                    alert.Type(typeName),
			Enabled:                 entry.Enabled,
			FailureThreshold:        entry.FailureThreshold,
			SuccessThreshold:        entry.SuccessThreshold,
			MinimumReminderInterval: minimumReminderInterval,
			Description:             entry.Description,
			SendOnResolved:          entry.SendOnResolved,
			ProviderOverride:        entry.ProviderOverride,
		}
		if err := parsed.ValidateAndSetDefaults(); err != nil {
			return nil, err
		}
		alerts = append(alerts, parsed)
	}
	return alerts, nil
}

func parseManagedAlertDuration(value string) (time.Duration, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) == 0 {
		return 0, nil
	}
	parsed, err := time.ParseDuration(trimmed)
	if err != nil {
		return 0, fmt.Errorf("invalid minimum reminder interval: %w", err)
	}
	return parsed, nil
}

func alertsToManagedPayload(alerts []*alert.Alert) []ManagedAlertPayload {
	serialized := make([]ManagedAlertPayload, 0, len(alerts))
	for _, entry := range alerts {
		if entry == nil {
			continue
		}
		alertPayload := ManagedAlertPayload{
			Type:             string(entry.Type),
			Enabled:          entry.Enabled,
			FailureThreshold: entry.FailureThreshold,
			SuccessThreshold: entry.SuccessThreshold,
			Description:      entry.Description,
			SendOnResolved:   entry.SendOnResolved,
		}
		if entry.MinimumReminderInterval > 0 {
			alertPayload.MinimumReminderInterval = entry.MinimumReminderInterval.String()
		}
		if len(entry.ProviderOverride) > 0 {
			alertPayload.ProviderOverride = make(map[string]any, len(entry.ProviderOverride))
			for key, value := range entry.ProviderOverride {
				alertPayload.ProviderOverride[key] = value
			}
		}
		serialized = append(serialized, alertPayload)
	}
	return serialized
}

func buildManagedAlertsFromTypes(types []string) ([]*alert.Alert, error) {
	payload := make([]ManagedAlertPayload, 0, len(types))
	for _, alertType := range types {
		trimmed := strings.TrimSpace(strings.ToLower(alertType))
		if len(trimmed) == 0 {
			continue
		}
		payload = append(payload, ManagedAlertPayload{
			Type: trimmed,
		})
	}
	return buildManagedAlertsFromPayload(payload)
}

func cloneManagedAlerts(alerts []*alert.Alert) []*alert.Alert {
	cloned := make([]*alert.Alert, 0, len(alerts))
	for _, entry := range alerts {
		if entry == nil {
			continue
		}
		alertClone := *entry
		if len(entry.ProviderOverride) > 0 {
			providerOverride := make(map[string]any, len(entry.ProviderOverride))
			for key, value := range entry.ProviderOverride {
				providerOverride[key] = value
			}
			alertClone.ProviderOverride = providerOverride
		}
		cloned = append(cloned, &alertClone)
	}
	return cloned
}
