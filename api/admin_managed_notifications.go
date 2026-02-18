package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/TwiN/gatus/v5/alerting"
	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/alerting/provider"
	"github.com/TwiN/gatus/v5/config"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

type ManagedNotification struct {
	Type                    string         `json:"type"`
	Configured              bool           `json:"configured"`
	UsedByEndpoints         int            `json:"usedByEndpoints"`
	UsedByExternalEndpoints int            `json:"usedByExternalEndpoints"`
	Config                  map[string]any `json:"config,omitempty"`
}

type ManagedNotificationListResponse struct {
	OverlayPath   string                `json:"overlayPath"`
	Notifications []ManagedNotification `json:"notifications"`
}

func GetManagedNotifications(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		notifications, err := buildManagedNotificationList(candidate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(ManagedNotificationListResponse{
			OverlayPath:   candidate.ManagedOverlayPath(),
			Notifications: notifications,
		})
	}
}

func PutManagedNotification(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notificationType := normalizeManagedNotificationType(c.Params("type"))
		fieldIndex, fieldType, err := getManagedNotificationField(notificationType)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		rawBody := c.Body()
		if len(rawBody) == 0 {
			rawBody = []byte("{}")
		}
		var payload map[string]any
		if err := json.Unmarshal(rawBody, &payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid payload: " + err.Error(),
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if candidate.Alerting == nil {
			candidate.Alerting = &alerting.Config{}
		}
		notificationProvider := reflect.New(fieldType.Elem())
		if err := yaml.Unmarshal(rawBody, notificationProvider.Interface()); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failed to decode provider configuration: " + err.Error(),
			})
		}
		alertProvider, ok := notificationProvider.Interface().(provider.AlertProvider)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "provider does not implement alerting provider interface",
			})
		}
		if err := alertProvider.Validate(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		reflect.ValueOf(candidate.Alerting).Elem().Field(fieldIndex).Set(notificationProvider)
		if err := validateManagedPayload(candidate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidateWithAlerting(cfg, candidate, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		serializedConfig, err := extractManagedNotificationConfig(notificationProvider)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		usedByEndpoints, usedByExternalEndpoints := countAlertTypeReferences(candidate, alert.Type(notificationType))
		return c.Status(fiber.StatusOK).JSON(ManagedNotification{
			Type:                    notificationType,
			Configured:              true,
			UsedByEndpoints:         usedByEndpoints,
			UsedByExternalEndpoints: usedByExternalEndpoints,
			Config:                  serializedConfig,
		})
	}
}

func DeleteManagedNotification(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notificationType := normalizeManagedNotificationType(c.Params("type"))
		fieldIndex, _, err := getManagedNotificationField(notificationType)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		usedByEndpoints, usedByExternalEndpoints := countAlertTypeReferences(candidate, alert.Type(notificationType))
		if usedByEndpoints > 0 || usedByExternalEndpoints > 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": fmt.Sprintf("notification type %s is still referenced by %d endpoint(s) and %d external endpoint(s)", notificationType, usedByEndpoints, usedByExternalEndpoints),
			})
		}
		if candidate.Alerting == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("notification type %s is not configured", notificationType),
			})
		}
		fieldValue := reflect.ValueOf(candidate.Alerting).Elem().Field(fieldIndex)
		if fieldValue.IsNil() {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("notification type %s is not configured", notificationType),
			})
		}
		fieldValue.Set(reflect.Zero(fieldValue.Type()))
		if isManagedAlertingConfigEmpty(candidate.Alerting) {
			candidate.Alerting = nil
		}
		if err := validateManagedPayload(candidate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidateWithAlerting(cfg, candidate, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func buildManagedNotificationList(candidate *config.Config) ([]ManagedNotification, error) {
	notifications := make([]ManagedNotification, 0)
	alertingType := reflect.TypeOf(alerting.Config{})
	var alertingValue reflect.Value
	if candidate.Alerting != nil {
		alertingValue = reflect.ValueOf(candidate.Alerting).Elem()
	}
	for i := 0; i < alertingType.NumField(); i++ {
		field := alertingType.Field(i)
		fieldTag := strings.Split(field.Tag.Get("yaml"), ",")[0]
		if len(fieldTag) == 0 || fieldTag == "-" {
			continue
		}
		usedByEndpoints, usedByExternalEndpoints := countAlertTypeReferences(candidate, alert.Type(fieldTag))
		notification := ManagedNotification{
			Type:                    fieldTag,
			UsedByEndpoints:         usedByEndpoints,
			UsedByExternalEndpoints: usedByExternalEndpoints,
		}
		if candidate.Alerting != nil {
			providerFieldValue := alertingValue.Field(i)
			if !providerFieldValue.IsNil() {
				serializedConfig, err := extractManagedNotificationConfig(providerFieldValue)
				if err != nil {
					return nil, err
				}
				notification.Configured = true
				notification.Config = serializedConfig
			}
		}
		notifications = append(notifications, notification)
	}
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].Type < notifications[j].Type
	})
	return notifications, nil
}

func extractManagedNotificationConfig(providerValue reflect.Value) (map[string]any, error) {
	serializedProvider, err := yaml.Marshal(providerValue.Interface())
	if err != nil {
		return nil, err
	}
	decoded := make(map[string]any)
	if err := yaml.Unmarshal(serializedProvider, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func getManagedNotificationField(notificationType string) (int, reflect.Type, error) {
	alertingType := reflect.TypeOf(alerting.Config{})
	for i := 0; i < alertingType.NumField(); i++ {
		field := alertingType.Field(i)
		fieldTag := strings.Split(field.Tag.Get("yaml"), ",")[0]
		if fieldTag == notificationType {
			return i, field.Type, nil
		}
	}
	return 0, nil, fmt.Errorf("unknown notification type %s", notificationType)
}

func normalizeManagedNotificationType(notificationType string) string {
	return strings.ToLower(strings.TrimSpace(notificationType))
}

func isManagedAlertingConfigEmpty(alertingConfig *alerting.Config) bool {
	if alertingConfig == nil {
		return true
	}
	value := reflect.ValueOf(alertingConfig).Elem()
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsNil() {
			return false
		}
	}
	return true
}

func countAlertTypeReferences(candidate *config.Config, notificationType alert.Type) (usedByEndpoints, usedByExternalEndpoints int) {
	for _, monitoredEndpoint := range candidate.Endpoints {
		for _, endpointAlert := range monitoredEndpoint.Alerts {
			if endpointAlert.Type == notificationType {
				usedByEndpoints++
			}
		}
	}
	for _, externalEndpoint := range candidate.ExternalEndpoints {
		for _, endpointAlert := range externalEndpoint.Alerts {
			if endpointAlert.Type == notificationType {
				usedByExternalEndpoints++
			}
		}
	}
	return usedByEndpoints, usedByExternalEndpoints
}
