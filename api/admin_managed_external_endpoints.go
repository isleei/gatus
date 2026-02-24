package api

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/endpoint/heartbeat"
	"github.com/gofiber/fiber/v2"
)

type ManagedExternalEndpointPayload struct {
	Enabled           *bool                 `json:"enabled,omitempty"`
	Name              string                `json:"name"`
	Group             string                `json:"group,omitempty"`
	Token             string                `json:"token,omitempty"`
	HeartbeatInterval string                `json:"heartbeatInterval,omitempty"`
	Alerts            []ManagedAlertPayload `json:"alerts,omitempty"`
}

type ManagedExternalEndpointResponse struct {
	Key               string                `json:"key"`
	Enabled           *bool                 `json:"enabled,omitempty"`
	Name              string                `json:"name"`
	Group             string                `json:"group,omitempty"`
	Token             string                `json:"token,omitempty"`
	HeartbeatInterval string                `json:"heartbeatInterval,omitempty"`
	Alerts            []ManagedAlertPayload `json:"alerts,omitempty"`
}

type ManagedExternalEndpointListResponse struct {
	OverlayPath       string                            `json:"overlayPath"`
	ExternalEndpoints []ManagedExternalEndpointResponse `json:"externalEndpoints"`
}

func GetManagedExternalEndpoints(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		externalEndpoints := make([]ManagedExternalEndpointResponse, 0, len(candidate.ExternalEndpoints))
		for _, externalEndpoint := range candidate.ExternalEndpoints {
			externalEndpoints = append(externalEndpoints, toManagedExternalEndpointResponse(externalEndpoint))
		}
		sort.Slice(externalEndpoints, func(i, j int) bool {
			return externalEndpoints[i].Key < externalEndpoints[j].Key
		})
		return c.Status(fiber.StatusOK).JSON(ManagedExternalEndpointListResponse{
			OverlayPath:       candidate.ManagedOverlayPath(),
			ExternalEndpoints: externalEndpoints,
		})
	}
}

func CreateManagedExternalEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload ManagedExternalEndpointPayload
		if err := c.BodyParser(&payload); err != nil {
			writeAdminAudit(c, cfg, "create", "external", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "create", "external", "", nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		externalEndpoint, err := buildManagedExternalEndpointFromPayload(&payload)
		if err != nil {
			writeAdminAudit(c, cfg, "create", "external", "", payload, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		candidate.ExternalEndpoints = append(candidate.ExternalEndpoints, externalEndpoint)
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "create", "external", externalEndpoint.Key(), payload, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "create", "external", externalEndpoint.Key(), payload, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		clearAdminDerivedCache()
		response := toManagedExternalEndpointResponse(externalEndpoint)
		writeAdminAudit(c, cfg, "create", "external", externalEndpoint.Key(), nil, response, nil)
		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func UpdateManagedExternalEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			err := errors.New("missing external endpoint key")
			writeAdminAudit(c, cfg, "update", "external", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var payload ManagedExternalEndpointPayload
		if err := c.BodyParser(&payload); err != nil {
			writeAdminAudit(c, cfg, "update", "external", targetKey, nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "update", "external", targetKey, nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		index := findExternalEndpointIndexByKey(candidate.ExternalEndpoints, targetKey)
		if index < 0 {
			err := fmt.Errorf("external endpoint with key %s not found", targetKey)
			writeAdminAudit(c, cfg, "update", "external", targetKey, payload, nil, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		before := toManagedExternalEndpointResponse(candidate.ExternalEndpoints[index])
		if err := applyManagedExternalEndpointPayload(candidate.ExternalEndpoints[index], &payload); err != nil {
			writeAdminAudit(c, cfg, "update", "external", targetKey, before, payload, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "update", "external", targetKey, before, payload, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "update", "external", targetKey, before, payload, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		clearAdminDerivedCache()
		response := toManagedExternalEndpointResponse(candidate.ExternalEndpoints[index])
		writeAdminAudit(c, cfg, "update", "external", response.Key, before, response, nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func DeleteManagedExternalEndpoint(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetKey := strings.TrimSpace(strings.ToLower(c.Params("key")))
		if len(targetKey) == 0 {
			err := errors.New("missing external endpoint key")
			writeAdminAudit(c, cfg, "delete", "external", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "delete", "external", targetKey, nil, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		index := findExternalEndpointIndexByKey(candidate.ExternalEndpoints, targetKey)
		if index < 0 {
			err := fmt.Errorf("external endpoint with key %s not found", targetKey)
			writeAdminAudit(c, cfg, "delete", "external", targetKey, nil, nil, err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		before := toManagedExternalEndpointResponse(candidate.ExternalEndpoints[index])
		candidate.ExternalEndpoints = append(candidate.ExternalEndpoints[:index], candidate.ExternalEndpoints[index+1:]...)
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "delete", "external", targetKey, before, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := persistManagedCandidate(cfg, candidate); err != nil {
			writeAdminAudit(c, cfg, "delete", "external", targetKey, before, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		clearAdminDerivedCache()
		writeAdminAudit(c, cfg, "delete", "external", targetKey, before, nil, nil)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func toManagedExternalEndpointResponse(externalEndpoint *endpoint.ExternalEndpoint) ManagedExternalEndpointResponse {
	response := ManagedExternalEndpointResponse{
		Key:     externalEndpoint.Key(),
		Enabled: externalEndpoint.Enabled,
		Name:    externalEndpoint.Name,
		Group:   externalEndpoint.Group,
		Token:   externalEndpoint.Token,
		Alerts:  alertsToManagedPayload(externalEndpoint.Alerts),
	}
	if externalEndpoint.Heartbeat.Interval > 0 {
		response.HeartbeatInterval = externalEndpoint.Heartbeat.Interval.String()
	}
	return response
}

func buildManagedExternalEndpointFromPayload(payload *ManagedExternalEndpointPayload) (*endpoint.ExternalEndpoint, error) {
	externalEndpoint := &endpoint.ExternalEndpoint{}
	if err := applyManagedExternalEndpointPayload(externalEndpoint, payload); err != nil {
		return nil, err
	}
	return externalEndpoint, nil
}

func applyManagedExternalEndpointPayload(externalEndpoint *endpoint.ExternalEndpoint, payload *ManagedExternalEndpointPayload) error {
	if externalEndpoint == nil || payload == nil {
		return errors.New("invalid external endpoint payload")
	}
	heartbeatInterval, err := parseManagedOptionalDuration(payload.HeartbeatInterval, "invalid heartbeat interval")
	if err != nil {
		return err
	}
	alerts, err := buildManagedAlertsFromPayload(payload.Alerts)
	if err != nil {
		return err
	}
	externalEndpoint.Enabled = payload.Enabled
	externalEndpoint.Name = strings.TrimSpace(payload.Name)
	externalEndpoint.Group = strings.TrimSpace(payload.Group)
	externalEndpoint.Token = strings.TrimSpace(payload.Token)
	externalEndpoint.Alerts = alerts
	externalEndpoint.Heartbeat = heartbeat.Config{
		Interval: heartbeatInterval,
	}
	return nil
}

func findExternalEndpointIndexByKey(externalEndpoints []*endpoint.ExternalEndpoint, targetKey string) int {
	for i := 0; i < len(externalEndpoints); i++ {
		if externalEndpoints[i].Key() == targetKey {
			return i
		}
	}
	return -1
}
