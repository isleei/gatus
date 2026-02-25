package api

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/client"
	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/gofiber/fiber/v2"
)

type ManagedEndpointPayload struct {
	Enabled     *bool                              `json:"enabled,omitempty"`
	Name        string                             `json:"name"`
	Group       string                             `json:"group,omitempty"`
	URL         string                             `json:"url"`
	Method      string                             `json:"method,omitempty"`
	Body        string                             `json:"body,omitempty"`
	GraphQL     bool                               `json:"graphql,omitempty"`
	Headers     map[string]string                  `json:"headers,omitempty"`
	Interval    string                             `json:"interval,omitempty"`
	Conditions  []string                           `json:"conditions"`
	Alerts      []ManagedAlertPayload              `json:"alerts,omitempty"`
	UI          *ManagedEndpointUIConfigPayload    `json:"ui,omitempty"`
	Client      *ManagedEndpointClientPayload      `json:"client,omitempty"`
	Certificate *ManagedEndpointCertificatePayload `json:"certificate,omitempty"`
	Tamper      *ManagedEndpointTamperPayload      `json:"tamper,omitempty"`
}

type ManagedEndpointResponse struct {
	Key         string                             `json:"key"`
	Type        endpoint.Type                      `json:"type"`
	Enabled     *bool                              `json:"enabled,omitempty"`
	Name        string                             `json:"name"`
	Group       string                             `json:"group,omitempty"`
	URL         string                             `json:"url"`
	Method      string                             `json:"method,omitempty"`
	Body        string                             `json:"body,omitempty"`
	GraphQL     bool                               `json:"graphql,omitempty"`
	Headers     map[string]string                  `json:"headers,omitempty"`
	Interval    string                             `json:"interval,omitempty"`
	Conditions  []string                           `json:"conditions"`
	Alerts      []ManagedAlertPayload              `json:"alerts,omitempty"`
	UI          *ManagedEndpointUIConfigPayload    `json:"ui,omitempty"`
	Client      *ManagedEndpointClientPayload      `json:"client,omitempty"`
	Certificate *ManagedEndpointCertificatePayload `json:"certificate,omitempty"`
	Tamper      *ManagedEndpointTamperPayload      `json:"tamper,omitempty"`
}

type ManagedEndpointClientPayload struct {
	Timeout        string `json:"timeout,omitempty"`
	Insecure       *bool  `json:"insecure,omitempty"`
	IgnoreRedirect *bool  `json:"ignoreRedirect,omitempty"`
}

type ManagedEndpointCertificatePayload struct {
	Enabled             bool   `json:"enabled"`
	ExpirationThreshold string `json:"expirationThreshold,omitempty"`
}

type ManagedEndpointTamperPayload struct {
	Enabled               bool     `json:"enabled"`
	BaselineSamples       int      `json:"baselineSamples,omitempty"`
	DriftThresholdPercent int64    `json:"driftThresholdPercent,omitempty"`
	ConsecutiveBreaches   int      `json:"consecutiveBreaches,omitempty"`
	RequiredSubstrings    []string `json:"requiredSubstrings,omitempty"`
	ForbiddenSubstrings   []string `json:"forbiddenSubstrings,omitempty"`
}

type ManagedEndpointListResponse struct {
	OverlayPath string                    `json:"overlayPath"`
	Endpoints   []ManagedEndpointResponse `json:"endpoints"`
}

func GetManagedEndpoints(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		endpoints := make([]ManagedEndpointResponse, 0, len(candidate.Endpoints))
		for _, monitoredEndpoint := range candidate.Endpoints {
			endpoints = append(endpoints, toManagedEndpointResponse(monitoredEndpoint))
		}
		sort.Slice(endpoints, func(i, j int) bool {
			return endpoints[i].Key < endpoints[j].Key
		})
		return c.Status(fiber.StatusOK).JSON(ManagedEndpointListResponse{
			OverlayPath: candidate.ManagedOverlayPath(),
			Endpoints:   endpoints,
		})
	}
}

func CreateManagedEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload ManagedEndpointPayload
		if err := c.BodyParser(&payload); err != nil {
			writeAdminAudit(c, cfg, "create", monitorEntityEndpoint, "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid payload: " + err.Error(),
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "create", monitorEntityEndpoint, "", payload, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		newEndpoint, err := buildManagedEndpointFromPayload(&payload)
		if err != nil {
			writeAdminAudit(c, cfg, "create", monitorEntityEndpoint, "", payload, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		candidate.Endpoints = append(candidate.Endpoints, newEndpoint)
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "create", monitorEntityEndpoint, newEndpoint.Key(), payload, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "create", monitorEntityEndpoint, newEndpoint.Key(), payload, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		clearAdminDerivedCache()
		response := toManagedEndpointResponse(newEndpoint)
		writeAdminAudit(c, cfg, "create", monitorEntityEndpoint, newEndpoint.Key(), nil, response, nil)
		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func UpdateManagedEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			err := errors.New("missing endpoint key")
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var payload ManagedEndpointPayload
		if err := c.BodyParser(&payload); err != nil {
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, targetKey, nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid payload: " + err.Error(),
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, targetKey, payload, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		index := findEndpointIndexByKey(candidate.Endpoints, targetKey)
		if index < 0 {
			err := fmt.Errorf("endpoint with key %s not found", targetKey)
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, targetKey, payload, nil, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		before := toManagedEndpointResponse(candidate.Endpoints[index])
		if err := applyManagedEndpointPayload(candidate.Endpoints[index], &payload); err != nil {
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, targetKey, before, payload, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, targetKey, before, payload, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, targetKey, before, payload, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		clearAdminDerivedCache()
		response := toManagedEndpointResponse(candidate.Endpoints[index])
		writeAdminAudit(c, cfg, "update", monitorEntityEndpoint, response.Key, before, response, nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func DeleteManagedEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			err := errors.New("missing endpoint key")
			writeAdminAudit(c, cfg, "delete", monitorEntityEndpoint, "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "delete", monitorEntityEndpoint, targetKey, nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		index := findEndpointIndexByKey(candidate.Endpoints, targetKey)
		if index < 0 {
			err := fmt.Errorf("endpoint with key %s not found", targetKey)
			writeAdminAudit(c, cfg, "delete", monitorEntityEndpoint, targetKey, nil, nil, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		before := toManagedEndpointResponse(candidate.Endpoints[index])
		candidate.Endpoints = append(candidate.Endpoints[:index], candidate.Endpoints[index+1:]...)
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "delete", monitorEntityEndpoint, targetKey, before, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "delete", monitorEntityEndpoint, targetKey, before, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		clearAdminDerivedCache()
		writeAdminAudit(c, cfg, "delete", monitorEntityEndpoint, targetKey, before, nil, nil)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func toManagedEndpointResponse(monitoredEndpoint *endpoint.Endpoint) ManagedEndpointResponse {
	response := ManagedEndpointResponse{
		Key:         monitoredEndpoint.Key(),
		Type:        monitoredEndpoint.Type(),
		Enabled:     monitoredEndpoint.Enabled,
		Name:        monitoredEndpoint.Name,
		Group:       monitoredEndpoint.Group,
		URL:         monitoredEndpoint.URL,
		Method:      monitoredEndpoint.Method,
		Body:        monitoredEndpoint.Body,
		GraphQL:     monitoredEndpoint.GraphQL,
		Headers:     make(map[string]string, len(monitoredEndpoint.Headers)),
		Interval:    monitoredEndpoint.Interval.String(),
		Conditions:  make([]string, 0, len(monitoredEndpoint.Conditions)),
		Alerts:      alertsToManagedPayload(monitoredEndpoint.Alerts),
		UI:          managedEndpointUIConfigToPayload(monitoredEndpoint.UIConfig),
		Client:      managedEndpointClientConfigToPayload(monitoredEndpoint.ClientConfig),
		Certificate: managedEndpointCertificateToPayload(monitoredEndpoint.Conditions),
		Tamper:      managedEndpointTamperToPayload(monitoredEndpoint.TamperConfig),
	}
	for headerName, headerValue := range monitoredEndpoint.Headers {
		response.Headers[headerName] = headerValue
	}
	for _, condition := range monitoredEndpoint.Conditions {
		response.Conditions = append(response.Conditions, string(condition))
	}
	return response
}

func buildManagedEndpointFromPayload(payload *ManagedEndpointPayload) (*endpoint.Endpoint, error) {
	monitoredEndpoint := &endpoint.Endpoint{}
	if err := applyManagedEndpointPayload(monitoredEndpoint, payload); err != nil {
		return nil, err
	}
	return monitoredEndpoint, nil
}

func applyManagedEndpointPayload(monitoredEndpoint *endpoint.Endpoint, payload *ManagedEndpointPayload) error {
	if monitoredEndpoint == nil || payload == nil {
		return errors.New("invalid endpoint payload")
	}
	interval, err := parseManagedEndpointInterval(payload.Interval)
	if err != nil {
		return err
	}
	conditions := make([]endpoint.Condition, 0, len(payload.Conditions))
	for _, condition := range payload.Conditions {
		trimmed := strings.TrimSpace(condition)
		if len(trimmed) == 0 {
			continue
		}
		conditions = append(conditions, endpoint.Condition(trimmed))
	}
	headers := make(map[string]string, len(payload.Headers))
	for headerName, headerValue := range payload.Headers {
		headers[headerName] = headerValue
	}
	monitoredEndpoint.Enabled = payload.Enabled
	monitoredEndpoint.Name = strings.TrimSpace(payload.Name)
	monitoredEndpoint.Group = strings.TrimSpace(payload.Group)
	monitoredEndpoint.URL = strings.TrimSpace(payload.URL)
	monitoredEndpoint.Method = strings.ToUpper(strings.TrimSpace(payload.Method))
	monitoredEndpoint.Body = payload.Body
	monitoredEndpoint.GraphQL = payload.GraphQL
	monitoredEndpoint.Headers = headers
	monitoredEndpoint.Interval = interval
	monitoredEndpoint.Conditions = conditions
	parsedAlerts, err := buildManagedAlertsFromPayload(payload.Alerts)
	if err != nil {
		return err
	}
	if payload.Certificate != nil {
		monitoredEndpoint.Conditions = removeManagedCertificateConditions(monitoredEndpoint.Conditions)
		if payload.Certificate.Enabled {
			threshold := strings.TrimSpace(payload.Certificate.ExpirationThreshold)
			if len(threshold) == 0 {
				return errors.New("certificate expiration threshold is required when certificate checks are enabled")
			}
			if _, err = time.ParseDuration(threshold); err != nil {
				if _, parseIntErr := strconv.ParseInt(threshold, 10, 64); parseIntErr != nil {
					return fmt.Errorf("invalid certificate expiration threshold: %w", err)
				}
			}
			monitoredEndpoint.Conditions = append(monitoredEndpoint.Conditions, endpoint.Condition(fmt.Sprintf("%s > %s", endpoint.CertificateExpirationPlaceholder, threshold)))
		}
	}
	if err = applyManagedEndpointClientPayload(monitoredEndpoint, payload.Client); err != nil {
		return err
	}
	if payload.Tamper != nil {
		monitoredEndpoint.TamperConfig = &endpoint.TamperConfig{
			Enabled:               payload.Tamper.Enabled,
			BaselineSamples:       payload.Tamper.BaselineSamples,
			DriftThresholdPercent: payload.Tamper.DriftThresholdPercent,
			ConsecutiveBreaches:   payload.Tamper.ConsecutiveBreaches,
			RequiredSubstrings:    append([]string(nil), payload.Tamper.RequiredSubstrings...),
			ForbiddenSubstrings:   append([]string(nil), payload.Tamper.ForbiddenSubstrings...),
		}
	}
	if payload.UI != nil {
		monitoredEndpoint.UIConfig = managedEndpointUIConfigFromPayload(payload.UI)
	}
	monitoredEndpoint.Alerts = parsedAlerts
	return nil
}

func parseManagedEndpointInterval(interval string) (time.Duration, error) {
	trimmed := strings.TrimSpace(interval)
	if len(trimmed) == 0 {
		return 0, nil
	}
	parsed, err := time.ParseDuration(trimmed)
	if err != nil {
		return 0, fmt.Errorf("invalid interval: %w", err)
	}
	return parsed, nil
}

func findEndpointIndexByKey(endpoints []*endpoint.Endpoint, targetKey string) int {
	for i := 0; i < len(endpoints); i++ {
		if endpoints[i].Key() == targetKey {
			return i
		}
	}
	return -1
}

func managedEndpointClientConfigToPayload(config *client.Config) *ManagedEndpointClientPayload {
	if config == nil {
		return nil
	}
	return &ManagedEndpointClientPayload{
		Timeout:        config.Timeout.String(),
		Insecure:       boolPointer(config.Insecure),
		IgnoreRedirect: boolPointer(config.IgnoreRedirect),
	}
}

func applyManagedEndpointClientPayload(monitoredEndpoint *endpoint.Endpoint, payload *ManagedEndpointClientPayload) error {
	if payload == nil {
		return nil
	}
	if monitoredEndpoint.ClientConfig == nil {
		monitoredEndpoint.ClientConfig = client.GetDefaultConfig()
	}
	if trimmedTimeout := strings.TrimSpace(payload.Timeout); len(trimmedTimeout) > 0 {
		timeout, err := time.ParseDuration(trimmedTimeout)
		if err != nil {
			return fmt.Errorf("invalid client timeout: %w", err)
		}
		monitoredEndpoint.ClientConfig.Timeout = timeout
	}
	if payload.Insecure != nil {
		monitoredEndpoint.ClientConfig.Insecure = *payload.Insecure
	}
	if payload.IgnoreRedirect != nil {
		monitoredEndpoint.ClientConfig.IgnoreRedirect = *payload.IgnoreRedirect
	}
	return nil
}

func managedEndpointTamperToPayload(config *endpoint.TamperConfig) *ManagedEndpointTamperPayload {
	if config == nil {
		return nil
	}
	return &ManagedEndpointTamperPayload{
		Enabled:               config.Enabled,
		BaselineSamples:       config.BaselineSamples,
		DriftThresholdPercent: config.DriftThresholdPercent,
		ConsecutiveBreaches:   config.ConsecutiveBreaches,
		RequiredSubstrings:    append([]string(nil), config.RequiredSubstrings...),
		ForbiddenSubstrings:   append([]string(nil), config.ForbiddenSubstrings...),
	}
}

func managedEndpointCertificateToPayload(conditions []endpoint.Condition) *ManagedEndpointCertificatePayload {
	threshold, enabled := extractManagedCertificateThreshold(conditions)
	return &ManagedEndpointCertificatePayload{
		Enabled:             enabled,
		ExpirationThreshold: threshold,
	}
}

func extractManagedCertificateThreshold(conditions []endpoint.Condition) (string, bool) {
	for _, condition := range conditions {
		conditionText := strings.TrimSpace(string(condition))
		parts := strings.Fields(conditionText)
		if len(parts) < 3 {
			continue
		}
		if !strings.EqualFold(parts[0], endpoint.CertificateExpirationPlaceholder) {
			continue
		}
		operator := parts[1]
		if operator != ">" && operator != ">=" {
			continue
		}
		threshold := strings.TrimSpace(strings.Join(parts[2:], " "))
		if len(threshold) == 0 {
			continue
		}
		return threshold, true
	}
	return "", false
}

func removeManagedCertificateConditions(conditions []endpoint.Condition) []endpoint.Condition {
	filtered := make([]endpoint.Condition, 0, len(conditions))
	for _, condition := range conditions {
		if strings.Contains(strings.ToUpper(string(condition)), endpoint.CertificateExpirationPlaceholder) {
			continue
		}
		filtered = append(filtered, condition)
	}
	return filtered
}

func boolPointer(value bool) *bool {
	return &value
}
