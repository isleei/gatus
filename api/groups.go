package api

import (
	"slices"
	"strings"

	"github.com/TwiN/gatus/v5/config"
	"github.com/gofiber/fiber/v2"
)

type GroupResponse struct {
	EndpointGroups []string `json:"endpointGroups"`
	SuiteGroups    []string `json:"suiteGroups"`
}

func appendUniqueGroup(groups []string, group string) []string {
	trimmed := strings.TrimSpace(group)
	if len(trimmed) == 0 || slices.Contains(groups, trimmed) {
		return groups
	}
	return append(groups, trimmed)
}

// Groups returns the list of configured endpoint and suite groups.
func Groups(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		response := GroupResponse{}
		for _, ep := range cfg.Endpoints {
			response.EndpointGroups = appendUniqueGroup(response.EndpointGroups, ep.Group)
		}
		for _, ep := range cfg.ExternalEndpoints {
			response.EndpointGroups = appendUniqueGroup(response.EndpointGroups, ep.Group)
		}
		for _, s := range cfg.Suites {
			response.SuiteGroups = appendUniqueGroup(response.SuiteGroups, s.Group)
		}
		if cfg.Security != nil && cfg.Security.Authorization != nil {
			if len(cfg.Security.Authorization.EndpointGroups) > 0 {
				response.EndpointGroups = cfg.Security.FilterAllowedEndpointGroups(response.EndpointGroups)
			}
			if len(cfg.Security.Authorization.SuiteGroups) > 0 {
				response.SuiteGroups = cfg.Security.FilterAllowedSuiteGroups(response.SuiteGroups)
			}
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
