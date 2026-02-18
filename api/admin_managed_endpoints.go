package api

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/gofiber/fiber/v2"
)

type ManagedEndpointPayload struct {
	Enabled    *bool             `json:"enabled,omitempty"`
	Name       string            `json:"name"`
	Group      string            `json:"group,omitempty"`
	URL        string            `json:"url"`
	Method     string            `json:"method,omitempty"`
	Body       string            `json:"body,omitempty"`
	GraphQL    bool              `json:"graphql,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Interval   string            `json:"interval,omitempty"`
	Conditions []string          `json:"conditions"`
}

type ManagedEndpointResponse struct {
	Key        string            `json:"key"`
	Type       endpoint.Type     `json:"type"`
	Enabled    *bool             `json:"enabled,omitempty"`
	Name       string            `json:"name"`
	Group      string            `json:"group,omitempty"`
	URL        string            `json:"url"`
	Method     string            `json:"method,omitempty"`
	Body       string            `json:"body,omitempty"`
	GraphQL    bool              `json:"graphql,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Interval   string            `json:"interval,omitempty"`
	Conditions []string          `json:"conditions"`
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
		newEndpoint, err := buildManagedEndpointFromPayload(&payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		candidate.Endpoints = append(candidate.Endpoints, newEndpoint)
		if err := validateManagedPayload(candidate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(toManagedEndpointResponse(newEndpoint))
	}
}

func UpdateManagedEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing endpoint key",
			})
		}
		var payload ManagedEndpointPayload
		if err := c.BodyParser(&payload); err != nil {
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
		index := findEndpointIndexByKey(candidate.Endpoints, targetKey)
		if index < 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("endpoint with key %s not found", targetKey),
			})
		}
		if err := applyManagedEndpointPayload(candidate.Endpoints[index], &payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := validateManagedPayload(candidate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(toManagedEndpointResponse(candidate.Endpoints[index]))
	}
}

func DeleteManagedEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing endpoint key",
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		index := findEndpointIndexByKey(candidate.Endpoints, targetKey)
		if index < 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("endpoint with key %s not found", targetKey),
			})
		}
		candidate.Endpoints = append(candidate.Endpoints[:index], candidate.Endpoints[index+1:]...)
		if err := validateManagedPayload(candidate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func toManagedEndpointResponse(monitoredEndpoint *endpoint.Endpoint) ManagedEndpointResponse {
	response := ManagedEndpointResponse{
		Key:        monitoredEndpoint.Key(),
		Type:       monitoredEndpoint.Type(),
		Enabled:    monitoredEndpoint.Enabled,
		Name:       monitoredEndpoint.Name,
		Group:      monitoredEndpoint.Group,
		URL:        monitoredEndpoint.URL,
		Method:     monitoredEndpoint.Method,
		Body:       monitoredEndpoint.Body,
		GraphQL:    monitoredEndpoint.GraphQL,
		Headers:    make(map[string]string, len(monitoredEndpoint.Headers)),
		Interval:   monitoredEndpoint.Interval.String(),
		Conditions: make([]string, 0, len(monitoredEndpoint.Conditions)),
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
