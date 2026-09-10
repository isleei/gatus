package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

type ManagedImportRequest struct {
	EntityType string               `json:"entityType" yaml:"entityType"`
	Mode       string               `json:"mode" yaml:"mode"` // merge or replace
	DryRun     bool                 `json:"dryRun" yaml:"dryRun"`
	Data       ManagedConfigPayload `json:"data" yaml:"data"`
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
		request, err := parseManagedImportRequest(c)
		if err != nil {
			writeAdminAudit(c, cfg, "import", "monitor", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: " + err.Error()})
		}
		// Query overrides help YAML overlay imports from the Admin UI
		if q := c.Query("entityType"); len(q) > 0 {
			request.EntityType = q
		}
		if q := c.Query("mode"); len(q) > 0 {
			request.Mode = q
		}
		if _, ok := c.Queries()["dryRun"]; ok {
			request.DryRun = c.QueryBool("dryRun")
		}
		request.EntityType = normalizeBatchEntityType(request.EntityType)
		request.Mode = normalizeImportMode(request.Mode)

		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			writeAdminAudit(c, cfg, "import", request.EntityType, "", request, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		preview, err := applyImportRequest(candidate, request)
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

func parseManagedImportRequest(c *fiber.Ctx) (*ManagedImportRequest, error) {
	body := c.Body()
	if len(body) == 0 {
		return nil, fmt.Errorf("empty body")
	}
	contentType := strings.ToLower(c.Get("Content-Type"))
	format := strings.ToLower(strings.TrimSpace(c.Query("format")))
	useYAML := strings.Contains(contentType, "yaml") || format == "yaml" || looksLikeYAMLImport(body)
	if useYAML {
		return parseManagedImportYAML(body)
	}
	var request ManagedImportRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, err
	}
	return &request, nil
}

func looksLikeYAMLImport(body []byte) bool {
	trimmed := strings.TrimSpace(string(body))
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return false
	}
	return strings.Contains(trimmed, "endpoints:") || strings.Contains(trimmed, "suites:") || strings.Contains(trimmed, "external-endpoints:") || strings.Contains(trimmed, "entityType:")
}

func parseManagedImportYAML(body []byte) (*ManagedImportRequest, error) {
	// Shape A: full ManagedImportRequest as YAML
	var wrapped ManagedImportRequest
	if err := yaml.Unmarshal(body, &wrapped); err == nil {
		hasData := len(wrapped.Data.Endpoints) > 0 || len(wrapped.Data.Suites) > 0 || len(wrapped.Data.ExternalEndpoints) > 0 || wrapped.Data.Alerting != nil
		if hasData || len(wrapped.EntityType) > 0 || len(wrapped.Mode) > 0 {
			return &wrapped, nil
		}
	}
	// Shape B: raw Gatus overlay fragment (endpoints / suites / external-endpoints / alerting)
	var payload ManagedConfigPayload
	if err := yaml.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("yaml: %w", err)
	}
	return &ManagedImportRequest{
		Mode: "merge",
		Data: payload,
	}, nil
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
