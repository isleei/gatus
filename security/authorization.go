package security

import (
	"slices"
	"strings"

	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
)

// AuthorizationConfig controls group-level access to endpoint/suite status APIs.
type AuthorizationConfig struct {
	// EndpointGroups is the list of endpoint groups that can be accessed.
	// If empty, all endpoint groups are allowed.
	EndpointGroups []string `yaml:"endpoint-groups,omitempty"`

	// SuiteGroups is the list of suite groups that can be accessed.
	// If empty, all suite groups are allowed.
	SuiteGroups []string `yaml:"suite-groups,omitempty"`
}

func normalizeGroup(group string) string {
	return strings.TrimSpace(strings.ToLower(group))
}

func containsGroup(allowedGroups []string, group string) bool {
	normalizedAllowedGroups := make([]string, 0, len(allowedGroups))
	for _, allowedGroup := range allowedGroups {
		normalizedAllowedGroups = append(normalizedAllowedGroups, normalizeGroup(allowedGroup))
	}
	return slices.Contains(normalizedAllowedGroups, normalizeGroup(group))
}

// IsEndpointGroupAllowed returns whether an endpoint group is authorized.
func (c *Config) IsEndpointGroupAllowed(group string) bool {
	if c == nil || c.Authorization == nil || len(c.Authorization.EndpointGroups) == 0 {
		return true
	}
	return containsGroup(c.Authorization.EndpointGroups, group)
}

// IsSuiteGroupAllowed returns whether a suite group is authorized.
func (c *Config) IsSuiteGroupAllowed(group string) bool {
	if c == nil || c.Authorization == nil || len(c.Authorization.SuiteGroups) == 0 {
		return true
	}
	return containsGroup(c.Authorization.SuiteGroups, group)
}

// FilterEndpointStatuses returns only statuses in allowed endpoint groups.
func (c *Config) FilterEndpointStatuses(statuses []*endpoint.Status) []*endpoint.Status {
	if c == nil || c.Authorization == nil || len(c.Authorization.EndpointGroups) == 0 {
		return statuses
	}
	filtered := make([]*endpoint.Status, 0, len(statuses))
	for _, status := range statuses {
		if c.IsEndpointGroupAllowed(status.Group) {
			filtered = append(filtered, status)
		}
	}
	return filtered
}

// FilterSuiteStatuses returns only statuses in allowed suite groups.
func (c *Config) FilterSuiteStatuses(statuses []*suite.Status) []*suite.Status {
	if c == nil || c.Authorization == nil || len(c.Authorization.SuiteGroups) == 0 {
		return statuses
	}
	filtered := make([]*suite.Status, 0, len(statuses))
	for _, status := range statuses {
		if c.IsSuiteGroupAllowed(status.Group) {
			filtered = append(filtered, status)
		}
	}
	return filtered
}

// FilterAllowedEndpointGroups filters a string slice by endpoint group authorization.
func (c *Config) FilterAllowedEndpointGroups(groups []string) []string {
	if c == nil || c.Authorization == nil || len(c.Authorization.EndpointGroups) == 0 {
		return groups
	}
	filtered := make([]string, 0, len(groups))
	for _, group := range groups {
		if c.IsEndpointGroupAllowed(group) {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

// FilterAllowedSuiteGroups filters a string slice by suite group authorization.
func (c *Config) FilterAllowedSuiteGroups(groups []string) []string {
	if c == nil || c.Authorization == nil || len(c.Authorization.SuiteGroups) == 0 {
		return groups
	}
	filtered := make([]string, 0, len(groups))
	for _, group := range groups {
		if c.IsSuiteGroupAllowed(group) {
			filtered = append(filtered, group)
		}
	}
	return filtered
}
