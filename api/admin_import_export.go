package api

import (
	"fmt"
	"strings"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/gofiber/fiber/v2"
)

type ManagedImportRequest struct {
	EntityType string               `json:"entityType"`
	Mode       string               `json:"mode"` // merge or replace
	DryRun     bool                 `json:"dryRun"`
	Data       ManagedConfigPayload `json:"data"`
}

type ImportPreviewResult struct {
	EntityType string         `json:"entityType"`
	Mode       string         `json:"mode"`
	DryRun     bool           `json:"dryRun"`
	Changes    map[string]int `json:"changes"`
	Message    string         `json:"message,omitempty"`
}

func ExportManagedConfiguration(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		entityType := normalizeBatchEntityType(c.Query("entityType"))
		response := ManagedConfigResponse{
			OverlayPath: candidate.ManagedOverlayPath(),
		}
		switch entityType {
		case monitorEntityEndpoint:
			response.Endpoints = candidate.Endpoints
		case monitorEntitySuite:
			response.Suites = candidate.Suites
		case monitorEntityExternal:
			response.ExternalEndpoints = candidate.ExternalEndpoints
		default:
			response.Alerting = candidate.Alerting
			response.Endpoints = candidate.Endpoints
			response.ExternalEndpoints = candidate.ExternalEndpoints
			response.Suites = candidate.Suites
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func ImportManagedConfiguration(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request ManagedImportRequest
		if err := c.BodyParser(&request); err != nil {
			writeAdminAudit(c, cfg, "import", "monitor", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		request.EntityType = normalizeBatchEntityType(request.EntityType)
		request.Mode = normalizeImportMode(request.Mode)

		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "import", request.EntityType, "", request, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		preview, err := applyImportRequest(candidate, &request)
		if err != nil {
			writeAdminAudit(c, cfg, "import", request.EntityType, "", request, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validateManagedPayload(candidate); err != nil {
			writeAdminAudit(c, cfg, "import", request.EntityType, "", request, preview, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "preview": preview})
		}
		if !request.DryRun {
			if err := persistManagedCandidateWithAlerting(cfg, candidate, true); err != nil {
				writeAdminAudit(c, cfg, "import", request.EntityType, "", request, preview, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error(), "preview": preview})
			}
			clearAdminDerivedCache()
			preview.Message = "Import applied."
		} else {
			preview.Message = "Dry run completed."
		}
		writeAdminAudit(c, cfg, "import", request.EntityType, "", request, preview, nil)
		return c.Status(fiber.StatusOK).JSON(preview)
	}
}

func applyImportRequest(candidate *config.Config, request *ManagedImportRequest) (*ImportPreviewResult, error) {
	if candidate == nil || request == nil {
		return nil, fmt.Errorf("invalid import request")
	}
	preview := &ImportPreviewResult{
		EntityType: request.EntityType,
		Mode:       request.Mode,
		DryRun:     request.DryRun,
		Changes:    make(map[string]int),
	}
	switch request.Mode {
	case "replace":
		applyImportReplace(candidate, request, preview)
	default:
		if err := applyImportMerge(candidate, request, preview); err != nil {
			return nil, err
		}
	}
	return preview, nil
}

func applyImportReplace(candidate *config.Config, request *ManagedImportRequest, preview *ImportPreviewResult) {
	switch request.EntityType {
	case monitorEntityEndpoint:
		preview.Changes["endpointsDeleted"] = len(candidate.Endpoints) - len(request.Data.Endpoints)
		if preview.Changes["endpointsDeleted"] < 0 {
			preview.Changes["endpointsDeleted"] = 0
		}
		candidate.Endpoints = request.Data.Endpoints
	case monitorEntitySuite:
		preview.Changes["suitesDeleted"] = len(candidate.Suites) - len(request.Data.Suites)
		if preview.Changes["suitesDeleted"] < 0 {
			preview.Changes["suitesDeleted"] = 0
		}
		candidate.Suites = request.Data.Suites
	case monitorEntityExternal:
		preview.Changes["externalDeleted"] = len(candidate.ExternalEndpoints) - len(request.Data.ExternalEndpoints)
		if preview.Changes["externalDeleted"] < 0 {
			preview.Changes["externalDeleted"] = 0
		}
		candidate.ExternalEndpoints = request.Data.ExternalEndpoints
	default:
		preview.Changes["endpointsDeleted"] = len(candidate.Endpoints)
		preview.Changes["suitesDeleted"] = len(candidate.Suites)
		preview.Changes["externalDeleted"] = len(candidate.ExternalEndpoints)
		candidate.Alerting = request.Data.Alerting
		candidate.Endpoints = request.Data.Endpoints
		candidate.Suites = request.Data.Suites
		candidate.ExternalEndpoints = request.Data.ExternalEndpoints
	}
}

func applyImportMerge(candidate *config.Config, request *ManagedImportRequest, preview *ImportPreviewResult) error {
	switch request.EntityType {
	case monitorEntityEndpoint:
		created, updated := mergeEndpoints(candidate, request.Data.Endpoints)
		preview.Changes["endpointsCreated"] = created
		preview.Changes["endpointsUpdated"] = updated
	case monitorEntitySuite:
		created, updated := mergeSuites(candidate, request.Data.Suites)
		preview.Changes["suitesCreated"] = created
		preview.Changes["suitesUpdated"] = updated
	case monitorEntityExternal:
		created, updated := mergeExternalEndpoints(candidate, request.Data.ExternalEndpoints)
		preview.Changes["externalCreated"] = created
		preview.Changes["externalUpdated"] = updated
	default:
		if request.Data.Alerting != nil {
			candidate.Alerting = request.Data.Alerting
		}
		endpointsCreated, endpointsUpdated := mergeEndpoints(candidate, request.Data.Endpoints)
		suitesCreated, suitesUpdated := mergeSuites(candidate, request.Data.Suites)
		externalCreated, externalUpdated := mergeExternalEndpoints(candidate, request.Data.ExternalEndpoints)
		preview.Changes["endpointsCreated"] = endpointsCreated
		preview.Changes["endpointsUpdated"] = endpointsUpdated
		preview.Changes["suitesCreated"] = suitesCreated
		preview.Changes["suitesUpdated"] = suitesUpdated
		preview.Changes["externalCreated"] = externalCreated
		preview.Changes["externalUpdated"] = externalUpdated
	}
	return nil
}

func mergeEndpoints(candidate *config.Config, incoming []*endpoint.Endpoint) (created, updated int) {
	for _, monitoredEndpoint := range incoming {
		if monitoredEndpoint == nil {
			continue
		}
		index := findEndpointIndexByKey(candidate.Endpoints, monitoredEndpoint.Key())
		if index < 0 {
			candidate.Endpoints = append(candidate.Endpoints, monitoredEndpoint)
			created++
			continue
		}
		candidate.Endpoints[index] = monitoredEndpoint
		updated++
	}
	return created, updated
}

func mergeSuites(candidate *config.Config, incoming []*suite.Suite) (created, updated int) {
	for _, monitoredSuite := range incoming {
		if monitoredSuite == nil {
			continue
		}
		index := findSuiteIndexByKey(candidate.Suites, monitoredSuite.Key())
		if index < 0 {
			candidate.Suites = append(candidate.Suites, monitoredSuite)
			created++
			continue
		}
		candidate.Suites[index] = monitoredSuite
		updated++
	}
	return created, updated
}

func mergeExternalEndpoints(candidate *config.Config, incoming []*endpoint.ExternalEndpoint) (created, updated int) {
	for _, externalEndpoint := range incoming {
		if externalEndpoint == nil {
			continue
		}
		index := findExternalEndpointIndexByKey(candidate.ExternalEndpoints, externalEndpoint.Key())
		if index < 0 {
			candidate.ExternalEndpoints = append(candidate.ExternalEndpoints, externalEndpoint)
			created++
			continue
		}
		candidate.ExternalEndpoints[index] = externalEndpoint
		updated++
	}
	return created, updated
}

func normalizeImportMode(mode string) string {
	if strings.EqualFold(strings.TrimSpace(mode), "replace") {
		return "replace"
	}
	return "merge"
}
