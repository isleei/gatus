package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/gofiber/fiber/v2"
)

type BatchOperationRequest struct {
	EntityType string                 `json:"entityType"`
	Keys       []string               `json:"keys"`
	Action     string                 `json:"action"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	DryRun     bool                   `json:"dryRun,omitempty"`
}

type BatchItemResult struct {
	Key     string `json:"key"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type BatchOperationResult struct {
	EntityType string            `json:"entityType"`
	Action     string            `json:"action"`
	DryRun     bool              `json:"dryRun"`
	Total      int               `json:"total"`
	Success    int               `json:"success"`
	Failed     int               `json:"failed"`
	Results    []BatchItemResult `json:"results"`
}

func ExecuteBatchOperation(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request BatchOperationRequest
		if err := c.BodyParser(&request); err != nil {
			writeAdminAudit(c, cfg, "batch", "monitor", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		request.EntityType = normalizeBatchEntityType(request.EntityType)
		request.Action = strings.ToLower(strings.TrimSpace(request.Action))
		if len(request.Keys) == 0 {
			err := errors.New("keys cannot be empty")
			writeAdminAudit(c, cfg, "batch", request.EntityType, "", request, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if !isBatchActionSupported(request.Action) {
			err := fmt.Errorf("unsupported action %s", request.Action)
			writeAdminAudit(c, cfg, "batch", request.EntityType, "", request, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "batch", request.EntityType, "", request, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		beforeSnapshot := map[string]interface{}{
			"entityType": request.EntityType,
			"action":     request.Action,
			"keys":       request.Keys,
		}

		results := make([]BatchItemResult, 0, len(request.Keys))
		for _, key := range request.Keys {
			key = strings.ToLower(strings.TrimSpace(key))
			if len(key) == 0 {
				continue
			}
			itemResult := BatchItemResult{
				Key:     key,
				Success: true,
			}
			if err := applyBatchOperationToKey(candidate, request.EntityType, key, request.Action, request.Payload); err != nil {
				itemResult.Success = false
				itemResult.Error = err.Error()
			}
			results = append(results, itemResult)
		}

		response := BatchOperationResult{
			EntityType: request.EntityType,
			Action:     request.Action,
			DryRun:     request.DryRun,
			Total:      len(results),
			Results:    results,
		}
		for _, result := range results {
			if result.Success {
				response.Success++
			} else {
				response.Failed++
			}
		}

		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "batch", request.EntityType, "", beforeSnapshot, response, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "result": response})
		}

		if !request.DryRun {
			if err := persistManagedCandidate(cfg, candidate); err != nil {
				writeAdminAudit(c, cfg, "batch", request.EntityType, "", beforeSnapshot, response, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "result": response})
			}
			clearAdminDerivedCache()
		}

		writeAdminAudit(c, cfg, "batch", request.EntityType, "", beforeSnapshot, response, nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func normalizeBatchEntityType(entityType string) string {
	switch strings.ToLower(strings.TrimSpace(entityType)) {
	case monitorEntityEndpoint, monitorEntitySuite, monitorEntityExternal:
		return strings.ToLower(strings.TrimSpace(entityType))
	default:
		return monitorEntityAll
	}
}

func isBatchActionSupported(action string) bool {
	switch action {
	case "enable", "disable", "delete", "set-group", "set-interval", "set-alert-types":
		return true
	default:
		return false
	}
}

func applyBatchOperationToKey(candidate *config.Config, entityType, key, action string, payload map[string]interface{}) error {
	handled := false
	if entityType == monitorEntityAll || entityType == monitorEntityEndpoint {
		index := findEndpointIndexByKey(candidate.Endpoints, key)
		if index >= 0 {
			handled = true
			if err := applyBatchOperationToEndpoint(candidate.Endpoints[index], action, payload); err != nil {
				return err
			}
			if action == "delete" {
				candidate.Endpoints = append(candidate.Endpoints[:index], candidate.Endpoints[index+1:]...)
			}
		}
	}
	if entityType == monitorEntityAll || entityType == monitorEntitySuite {
		index := findSuiteIndexByKey(candidate.Suites, key)
		if index >= 0 {
			handled = true
			if err := applyBatchOperationToSuite(candidate.Suites[index], action, payload); err != nil {
				return err
			}
			if action == "delete" {
				candidate.Suites = append(candidate.Suites[:index], candidate.Suites[index+1:]...)
			}
		}
	}
	if entityType == monitorEntityAll || entityType == monitorEntityExternal {
		index := findExternalEndpointIndexByKey(candidate.ExternalEndpoints, key)
		if index >= 0 {
			handled = true
			if err := applyBatchOperationToExternalEndpoint(candidate.ExternalEndpoints[index], action, payload); err != nil {
				return err
			}
			if action == "delete" {
				candidate.ExternalEndpoints = append(candidate.ExternalEndpoints[:index], candidate.ExternalEndpoints[index+1:]...)
			}
		}
	}
	if !handled {
		return fmt.Errorf("monitor with key %s not found", key)
	}
	return nil
}

func applyBatchOperationToEndpoint(monitoredEndpoint *endpoint.Endpoint, action string, payload map[string]interface{}) error {
	switch action {
	case "enable":
		enabled := true
		monitoredEndpoint.Enabled = &enabled
	case "disable":
		enabled := false
		monitoredEndpoint.Enabled = &enabled
	case "delete":
		return nil
	case "set-group":
		group, err := parseBatchStringPayload(payload, "group")
		if err != nil {
			return err
		}
		monitoredEndpoint.Group = group
	case "set-interval":
		interval, err := parseBatchDurationPayload(payload, "interval")
		if err != nil {
			return err
		}
		monitoredEndpoint.Interval = interval
	case "set-alert-types":
		alertTypes, err := parseBatchAlertTypes(payload)
		if err != nil {
			return err
		}
		alerts, err := buildManagedAlertsFromTypes(alertTypes)
		if err != nil {
			return err
		}
		monitoredEndpoint.Alerts = alerts
	default:
		return fmt.Errorf("unsupported action %s", action)
	}
	return nil
}

func applyBatchOperationToSuite(monitoredSuite *suite.Suite, action string, payload map[string]interface{}) error {
	switch action {
	case "enable":
		enabled := true
		monitoredSuite.Enabled = &enabled
	case "disable":
		enabled := false
		monitoredSuite.Enabled = &enabled
	case "delete":
		return nil
	case "set-group":
		group, err := parseBatchStringPayload(payload, "group")
		if err != nil {
			return err
		}
		monitoredSuite.Group = group
		for _, endpointConfig := range monitoredSuite.Endpoints {
			endpointConfig.Group = group
		}
	case "set-interval":
		interval, err := parseBatchDurationPayload(payload, "interval")
		if err != nil {
			return err
		}
		monitoredSuite.Interval = interval
	case "set-alert-types":
		alertTypes, err := parseBatchAlertTypes(payload)
		if err != nil {
			return err
		}
		alerts, err := buildManagedAlertsFromTypes(alertTypes)
		if err != nil {
			return err
		}
		for _, endpointConfig := range monitoredSuite.Endpoints {
			endpointConfig.Alerts = cloneManagedAlerts(alerts)
		}
	default:
		return fmt.Errorf("unsupported action %s", action)
	}
	return nil
}

func applyBatchOperationToExternalEndpoint(externalEndpoint *endpoint.ExternalEndpoint, action string, payload map[string]interface{}) error {
	switch action {
	case "enable":
		enabled := true
		externalEndpoint.Enabled = &enabled
	case "disable":
		enabled := false
		externalEndpoint.Enabled = &enabled
	case "delete":
		return nil
	case "set-group":
		group, err := parseBatchStringPayload(payload, "group")
		if err != nil {
			return err
		}
		externalEndpoint.Group = group
	case "set-interval":
		interval, err := parseBatchDurationPayload(payload, "interval")
		if err != nil {
			return err
		}
		externalEndpoint.Heartbeat.Interval = interval
	case "set-alert-types":
		alertTypes, err := parseBatchAlertTypes(payload)
		if err != nil {
			return err
		}
		alerts, err := buildManagedAlertsFromTypes(alertTypes)
		if err != nil {
			return err
		}
		externalEndpoint.Alerts = alerts
	default:
		return fmt.Errorf("unsupported action %s", action)
	}
	return nil
}

func parseBatchStringPayload(payload map[string]interface{}, field string) (string, error) {
	if payload == nil {
		return "", fmt.Errorf("missing payload field %s", field)
	}
	raw, exists := payload[field]
	if !exists {
		return "", fmt.Errorf("missing payload field %s", field)
	}
	parsed, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("payload field %s must be a string", field)
	}
	return strings.TrimSpace(parsed), nil
}

func parseBatchDurationPayload(payload map[string]interface{}, field string) (time.Duration, error) {
	value, err := parseBatchStringPayload(payload, field)
	if err != nil {
		return 0, err
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("payload field %s has invalid duration: %w", field, err)
	}
	return parsed, nil
}

func parseBatchAlertTypes(payload map[string]interface{}) ([]string, error) {
	if payload == nil {
		return nil, errors.New("missing payload field alertTypes")
	}
	raw, exists := payload["alertTypes"]
	if !exists {
		return nil, errors.New("missing payload field alertTypes")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, errors.New("payload field alertTypes must be an array")
	}
	result := make([]string, 0, len(list))
	for _, item := range list {
		value, ok := item.(string)
		if !ok {
			return nil, errors.New("payload field alertTypes must contain strings")
		}
		value = strings.TrimSpace(strings.ToLower(value))
		if value == "" {
			continue
		}
		result = append(result, value)
	}
	return result, nil
}
