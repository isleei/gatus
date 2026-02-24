package api

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/gofiber/fiber/v2"
)

type ManagedSuiteEndpointPayload struct {
	Enabled    *bool                 `json:"enabled,omitempty"`
	Name       string                `json:"name"`
	URL        string                `json:"url"`
	Method     string                `json:"method,omitempty"`
	Body       string                `json:"body,omitempty"`
	GraphQL    bool                  `json:"graphql,omitempty"`
	Headers    map[string]string     `json:"headers,omitempty"`
	Interval   string                `json:"interval,omitempty"`
	Conditions []string              `json:"conditions"`
	Store      map[string]string     `json:"store,omitempty"`
	AlwaysRun  bool                  `json:"alwaysRun,omitempty"`
	Alerts     []ManagedAlertPayload `json:"alerts,omitempty"`
}

type ManagedSuitePayload struct {
	Enabled   *bool                         `json:"enabled,omitempty"`
	Name      string                        `json:"name"`
	Group     string                        `json:"group,omitempty"`
	Interval  string                        `json:"interval,omitempty"`
	Timeout   string                        `json:"timeout,omitempty"`
	Context   map[string]interface{}        `json:"context,omitempty"`
	Endpoints []ManagedSuiteEndpointPayload `json:"endpoints"`
}

type ManagedSuiteResponse struct {
	Key       string                        `json:"key"`
	Enabled   *bool                         `json:"enabled,omitempty"`
	Name      string                        `json:"name"`
	Group     string                        `json:"group,omitempty"`
	Interval  string                        `json:"interval,omitempty"`
	Timeout   string                        `json:"timeout,omitempty"`
	Context   map[string]interface{}        `json:"context,omitempty"`
	Endpoints []ManagedSuiteEndpointPayload `json:"endpoints"`
}

type ManagedSuiteListResponse struct {
	OverlayPath string                 `json:"overlayPath"`
	Suites      []ManagedSuiteResponse `json:"suites"`
}

func GetManagedSuites(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		suites := make([]ManagedSuiteResponse, 0, len(candidate.Suites))
		for _, monitoredSuite := range candidate.Suites {
			suites = append(suites, toManagedSuiteResponse(monitoredSuite))
		}
		sort.Slice(suites, func(i, j int) bool {
			return suites[i].Key < suites[j].Key
		})
		return c.Status(fiber.StatusOK).JSON(ManagedSuiteListResponse{
			OverlayPath: candidate.ManagedOverlayPath(),
			Suites:      suites,
		})
	}
}

func CreateManagedSuite(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload ManagedSuitePayload
		if err := c.BodyParser(&payload); err != nil {
			writeAdminAudit(c, cfg, "create", "suite", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "create", "suite", "", nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		newSuite, err := buildManagedSuiteFromPayload(&payload)
		if err != nil {
			writeAdminAudit(c, cfg, "create", "suite", "", payload, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		candidate.Suites = append(candidate.Suites, newSuite)
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "create", "suite", newSuite.Key(), payload, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "create", "suite", newSuite.Key(), payload, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		clearAdminDerivedCache()
		response := toManagedSuiteResponse(newSuite)
		writeAdminAudit(c, cfg, "create", "suite", newSuite.Key(), nil, response, nil)
		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func UpdateManagedSuite(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			err := errors.New("missing suite key")
			writeAdminAudit(c, cfg, "update", "suite", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var payload ManagedSuitePayload
		if err := c.BodyParser(&payload); err != nil {
			writeAdminAudit(c, cfg, "update", "suite", targetKey, nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "update", "suite", targetKey, nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		index := findSuiteIndexByKey(candidate.Suites, targetKey)
		if index < 0 {
			err := fmt.Errorf("suite with key %s not found", targetKey)
			writeAdminAudit(c, cfg, "update", "suite", targetKey, payload, nil, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		before := toManagedSuiteResponse(candidate.Suites[index])
		if err := applyManagedSuitePayload(candidate.Suites[index], &payload); err != nil {
			writeAdminAudit(c, cfg, "update", "suite", targetKey, before, payload, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "update", "suite", targetKey, before, payload, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "update", "suite", targetKey, before, payload, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		clearAdminDerivedCache()
		response := toManagedSuiteResponse(candidate.Suites[index])
		writeAdminAudit(c, cfg, "update", "suite", response.Key, before, response, nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func DeleteManagedSuite(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			err := errors.New("missing suite key")
			writeAdminAudit(c, cfg, "delete", "suite", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "delete", "suite", targetKey, nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		index := findSuiteIndexByKey(candidate.Suites, targetKey)
		if index < 0 {
			err := fmt.Errorf("suite with key %s not found", targetKey)
			writeAdminAudit(c, cfg, "delete", "suite", targetKey, nil, nil, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		before := toManagedSuiteResponse(candidate.Suites[index])
		candidate.Suites = append(candidate.Suites[:index], candidate.Suites[index+1:]...)
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "delete", "suite", targetKey, before, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "delete", "suite", targetKey, before, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		clearAdminDerivedCache()
		writeAdminAudit(c, cfg, "delete", "suite", targetKey, before, nil, nil)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func toManagedSuiteResponse(monitoredSuite *suite.Suite) ManagedSuiteResponse {
	response := ManagedSuiteResponse{
		Key:       monitoredSuite.Key(),
		Enabled:   monitoredSuite.Enabled,
		Name:      monitoredSuite.Name,
		Group:     monitoredSuite.Group,
		Interval:  monitoredSuite.Interval.String(),
		Timeout:   monitoredSuite.Timeout.String(),
		Context:   make(map[string]interface{}, len(monitoredSuite.InitialContext)),
		Endpoints: make([]ManagedSuiteEndpointPayload, 0, len(monitoredSuite.Endpoints)),
	}
	for key, value := range monitoredSuite.InitialContext {
		response.Context[key] = value
	}
	for _, endpointConfig := range monitoredSuite.Endpoints {
		conditions := make([]string, 0, len(endpointConfig.Conditions))
		for _, condition := range endpointConfig.Conditions {
			conditions = append(conditions, string(condition))
		}
		headers := make(map[string]string, len(endpointConfig.Headers))
		for headerName, headerValue := range endpointConfig.Headers {
			headers[headerName] = headerValue
		}
		storeMap := make(map[string]string, len(endpointConfig.Store))
		for storeKey, storeValue := range endpointConfig.Store {
			storeMap[storeKey] = storeValue
		}
		response.Endpoints = append(response.Endpoints, ManagedSuiteEndpointPayload{
			Enabled:    endpointConfig.Enabled,
			Name:       endpointConfig.Name,
			URL:        endpointConfig.URL,
			Method:     endpointConfig.Method,
			Body:       endpointConfig.Body,
			GraphQL:    endpointConfig.GraphQL,
			Headers:    headers,
			Interval:   endpointConfig.Interval.String(),
			Conditions: conditions,
			Store:      storeMap,
			AlwaysRun:  endpointConfig.AlwaysRun,
			Alerts:     alertsToManagedPayload(endpointConfig.Alerts),
		})
	}
	return response
}

func buildManagedSuiteFromPayload(payload *ManagedSuitePayload) (*suite.Suite, error) {
	monitoredSuite := &suite.Suite{}
	if err := applyManagedSuitePayload(monitoredSuite, payload); err != nil {
		return nil, err
	}
	return monitoredSuite, nil
}

func applyManagedSuitePayload(monitoredSuite *suite.Suite, payload *ManagedSuitePayload) error {
	if monitoredSuite == nil || payload == nil {
		return errors.New("invalid suite payload")
	}
	interval, err := parseManagedOptionalDuration(payload.Interval, "invalid suite interval")
	if err != nil {
		return err
	}
	timeout, err := parseManagedOptionalDuration(payload.Timeout, "invalid suite timeout")
	if err != nil {
		return err
	}
	endpoints := make([]*endpoint.Endpoint, 0, len(payload.Endpoints))
	for _, endpointPayload := range payload.Endpoints {
		endpointInterval, err := parseManagedOptionalDuration(endpointPayload.Interval, "invalid suite endpoint interval")
		if err != nil {
			return err
		}
		conditions := make([]endpoint.Condition, 0, len(endpointPayload.Conditions))
		for _, condition := range endpointPayload.Conditions {
			trimmed := strings.TrimSpace(condition)
			if len(trimmed) == 0 {
				continue
			}
			conditions = append(conditions, endpoint.Condition(trimmed))
		}
		headers := make(map[string]string, len(endpointPayload.Headers))
		for headerName, headerValue := range endpointPayload.Headers {
			headers[headerName] = headerValue
		}
		storeMap := make(map[string]string, len(endpointPayload.Store))
		for storeKey, storeValue := range endpointPayload.Store {
			storeMap[storeKey] = storeValue
		}
		alerts, err := buildManagedAlertsFromPayload(endpointPayload.Alerts)
		if err != nil {
			return err
		}
		endpoints = append(endpoints, &endpoint.Endpoint{
			Enabled:    endpointPayload.Enabled,
			Name:       strings.TrimSpace(endpointPayload.Name),
			URL:        strings.TrimSpace(endpointPayload.URL),
			Method:     strings.ToUpper(strings.TrimSpace(endpointPayload.Method)),
			Body:       endpointPayload.Body,
			GraphQL:    endpointPayload.GraphQL,
			Headers:    headers,
			Interval:   endpointInterval,
			Conditions: conditions,
			Store:      storeMap,
			AlwaysRun:  endpointPayload.AlwaysRun,
			Alerts:     alerts,
		})
	}
	contextMap := make(map[string]interface{}, len(payload.Context))
	for key, value := range payload.Context {
		contextMap[key] = value
	}
	monitoredSuite.Enabled = payload.Enabled
	monitoredSuite.Name = strings.TrimSpace(payload.Name)
	monitoredSuite.Group = strings.TrimSpace(payload.Group)
	monitoredSuite.Interval = interval
	monitoredSuite.Timeout = timeout
	monitoredSuite.InitialContext = contextMap
	monitoredSuite.Endpoints = endpoints
	return nil
}

func findSuiteIndexByKey(suites []*suite.Suite, targetKey string) int {
	for i := 0; i < len(suites); i++ {
		if suites[i].Key() == targetKey {
			return i
		}
	}
	return -1
}

func parseManagedOptionalDuration(value, errorPrefix string) (time.Duration, error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) == 0 {
		return 0, nil
	}
	parsed, err := time.ParseDuration(trimmed)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errorPrefix, err)
	}
	return parsed, nil
}
