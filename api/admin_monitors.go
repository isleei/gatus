package api

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common/paging"
	"github.com/gofiber/fiber/v2"
)

const (
	monitorEntityAll      = "all"
	monitorEntityEndpoint = "endpoint"
	monitorEntitySuite    = "suite"
	monitorEntityExternal = "external"
)

type MonitorSummary struct {
	Key               string     `json:"key"`
	EntityType        string     `json:"entityType"`
	Type              string     `json:"type"`
	Name              string     `json:"name"`
	Group             string     `json:"group,omitempty"`
	Enabled           bool       `json:"enabled"`
	Status            string     `json:"status"`
	URL               string     `json:"url,omitempty"`
	Steps             int        `json:"steps,omitempty"`
	Interval          string     `json:"interval,omitempty"`
	LastCheck         *time.Time `json:"lastCheck,omitempty"`
	Duration          string     `json:"duration,omitempty"`
	DurationMs        int64      `json:"durationMs,omitempty"`
	NotificationTypes []string   `json:"notificationTypes,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

type MonitorKPI struct {
	Total     int `json:"total"`
	Unhealthy int `json:"unhealthy"`
	Disabled  int `json:"disabled"`
	Unknown   int `json:"unknown"`
}

type MonitorListResponse struct {
	Items    []MonitorSummary `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	KPI      MonitorKPI       `json:"kpi"`
	Groups   []string         `json:"groups"`
}

func GetAdminMonitors(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cacheKey := buildAdminCacheKey("monitors", c)
		if value, exists := cache.Get(cacheKey); exists {
			c.Set("Content-Type", "application/json")
			return c.Status(fiber.StatusOK).Send(value.([]byte))
		}

		response, err := buildMonitorListResponse(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		serialized, err := json.Marshal(response)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		cache.SetWithTTL(cacheKey, serialized, adminListCacheTTL)
		c.Set("Content-Type", "application/json")
		return c.Status(fiber.StatusOK).Send(serialized)
	}
}

func buildMonitorListResponse(cfg *config.Config, c *fiber.Ctx) (MonitorListResponse, error) {
	candidate, err := loadManagedCandidate(cfg)
	if err != nil {
		return MonitorListResponse{}, err
	}

	entityType := normalizeMonitorEntityType(c.Query("entityType"))
	query := strings.ToLower(strings.TrimSpace(c.Query("q")))
	groupFilter := strings.TrimSpace(c.Query("group"))
	enabledFilter := normalizeMonitorEnabledFilter(c.Query("enabled"))
	statusFilter := normalizeMonitorStatusFilter(c.Query("status"))
	sortBy := normalizeMonitorSortBy(c.Query("sortBy"))
	sortDir := normalizeMonitorSortDirection(c.Query("sortDir"))
	page, pageSize := parseMonitorPagination(c.Query("page"), c.Query("pageSize"))

	endpointStatuses, err := store.Get().GetAllEndpointStatuses(paging.NewEndpointStatusParams().WithResults(1, 1))
	if err != nil {
		return MonitorListResponse{}, err
	}
	suiteStatuses, err := store.Get().GetAllSuiteStatuses(paging.NewSuiteStatusParams().WithPagination(1, 1))
	if err != nil {
		return MonitorListResponse{}, err
	}
	endpointLatestResults := make(map[string]*endpoint.Result, len(endpointStatuses))
	for _, status := range endpointStatuses {
		if status == nil || len(status.Results) == 0 {
			continue
		}
		endpointLatestResults[status.Key] = status.Results[len(status.Results)-1]
	}
	suiteLatestResults := make(map[string]*suite.Result, len(suiteStatuses))
	for _, status := range suiteStatuses {
		if status == nil || len(status.Results) == 0 {
			continue
		}
		suiteLatestResults[status.Key] = status.Results[len(status.Results)-1]
	}

	monitors := make([]MonitorSummary, 0, len(candidate.Endpoints)+len(candidate.Suites)+len(candidate.ExternalEndpoints))
	for _, monitoredEndpoint := range candidate.Endpoints {
		latestResult := endpointLatestResults[monitoredEndpoint.Key()]
		summary := MonitorSummary{
			Key:               monitoredEndpoint.Key(),
			EntityType:        monitorEntityEndpoint,
			Type:              string(monitoredEndpoint.Type()),
			Name:              monitoredEndpoint.Name,
			Group:             monitoredEndpoint.Group,
			Enabled:           monitoredEndpoint.IsEnabled(),
			URL:               monitoredEndpoint.URL,
			Interval:          monitoredEndpoint.Interval.String(),
			NotificationTypes: extractAlertTypes(monitoredEndpoint.Alerts),
		}
		enrichMonitorSummaryWithEndpointResult(&summary, latestResult)
		monitors = append(monitors, summary)
	}
	for _, monitoredSuite := range candidate.Suites {
		latestResult := suiteLatestResults[monitoredSuite.Key()]
		summary := MonitorSummary{
			Key:               monitoredSuite.Key(),
			EntityType:        monitorEntitySuite,
			Type:              "SUITE",
			Name:              monitoredSuite.Name,
			Group:             monitoredSuite.Group,
			Enabled:           monitoredSuite.IsEnabled(),
			Steps:             len(monitoredSuite.Endpoints),
			Interval:          monitoredSuite.Interval.String(),
			NotificationTypes: extractSuiteAlertTypes(monitoredSuite.Endpoints),
		}
		enrichMonitorSummaryWithSuiteResult(&summary, latestResult)
		monitors = append(monitors, summary)
	}
	for _, externalEndpoint := range candidate.ExternalEndpoints {
		latestResult := endpointLatestResults[externalEndpoint.Key()]
		heartbeatInterval := ""
		if externalEndpoint.Heartbeat.Interval > 0 {
			heartbeatInterval = externalEndpoint.Heartbeat.Interval.String()
		}
		summary := MonitorSummary{
			Key:               externalEndpoint.Key(),
			EntityType:        monitorEntityExternal,
			Type:              "EXTERNAL",
			Name:              externalEndpoint.Name,
			Group:             externalEndpoint.Group,
			Enabled:           externalEndpoint.IsEnabled(),
			URL:               fmt.Sprintf("external://%s", externalEndpoint.Key()),
			Interval:          heartbeatInterval,
			NotificationTypes: extractAlertTypes(externalEndpoint.Alerts),
		}
		enrichMonitorSummaryWithEndpointResult(&summary, latestResult)
		monitors = append(monitors, summary)
	}

	filtered := filterMonitors(monitors, entityType, query, groupFilter, enabledFilter, statusFilter)
	sortMonitors(filtered, sortBy, sortDir)
	kpi := buildMonitorKPI(filtered)
	groups := collectMonitorGroups(filtered)

	total := len(filtered)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	items := filtered[start:end]

	return MonitorListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		KPI:      kpi,
		Groups:   groups,
	}, nil
}

func enrichMonitorSummaryWithEndpointResult(summary *MonitorSummary, result *endpoint.Result) {
	if summary == nil {
		return
	}
	if !summary.Enabled {
		summary.Status = "disabled"
		return
	}
	if result == nil {
		summary.Status = "unknown"
		return
	}
	if result.Success {
		summary.Status = "healthy"
	} else {
		summary.Status = "unhealthy"
	}
	timestamp := result.Timestamp
	summary.LastCheck = &timestamp
	summary.UpdatedAt = &timestamp
	summary.Duration = result.Duration.String()
	summary.DurationMs = result.Duration.Milliseconds()
}

func enrichMonitorSummaryWithSuiteResult(summary *MonitorSummary, result *suite.Result) {
	if summary == nil {
		return
	}
	if !summary.Enabled {
		summary.Status = "disabled"
		return
	}
	if result == nil {
		summary.Status = "unknown"
		return
	}
	if result.Success {
		summary.Status = "healthy"
	} else {
		summary.Status = "unhealthy"
	}
	timestamp := result.Timestamp
	summary.LastCheck = &timestamp
	summary.UpdatedAt = &timestamp
	summary.Duration = result.Duration.String()
	summary.DurationMs = result.Duration.Milliseconds()
}

func filterMonitors(monitors []MonitorSummary, entityType, query, groupFilter, enabledFilter, statusFilter string) []MonitorSummary {
	filtered := make([]MonitorSummary, 0, len(monitors))
	normalizedGroup := normalizeMonitorGroupFilter(groupFilter)
	for _, monitor := range monitors {
		if entityType != monitorEntityAll && monitor.EntityType != entityType {
			continue
		}
		if enabledFilter == "true" && !monitor.Enabled {
			continue
		}
		if enabledFilter == "false" && monitor.Enabled {
			continue
		}
		if statusFilter != "all" && monitor.Status != statusFilter {
			continue
		}
		if normalizedGroup != "" && strings.ToLower(strings.TrimSpace(monitor.Group)) != normalizedGroup {
			continue
		}
		if query != "" {
			haystack := strings.ToLower(strings.Join([]string{
				monitor.Key,
				monitor.Name,
				monitor.Group,
				monitor.Type,
				monitor.URL,
				strings.Join(monitor.NotificationTypes, " "),
			}, " "))
			if !strings.Contains(haystack, query) {
				continue
			}
		}
		filtered = append(filtered, monitor)
	}
	return filtered
}

func sortMonitors(monitors []MonitorSummary, sortBy, sortDir string) {
	ascending := sortDir == "asc"
	statusWeight := func(status string) int {
		switch status {
		case "unhealthy":
			return 0
		case "unknown":
			return 1
		case "healthy":
			return 2
		case "disabled":
			return 3
		default:
			return 4
		}
	}
	compareTime := func(left, right *time.Time) int {
		if left == nil && right == nil {
			return 0
		}
		if left == nil {
			return -1
		}
		if right == nil {
			return 1
		}
		if left.Before(*right) {
			return -1
		}
		if left.After(*right) {
			return 1
		}
		return 0
	}
	sort.Slice(monitors, func(i, j int) bool {
		left := monitors[i]
		right := monitors[j]
		order := 0
		switch sortBy {
		case "name":
			order = strings.Compare(strings.ToLower(left.Name), strings.ToLower(right.Name))
		case "group":
			order = strings.Compare(strings.ToLower(left.Group), strings.ToLower(right.Group))
		case "status":
			order = statusWeight(left.Status) - statusWeight(right.Status)
		case "interval":
			order = strings.Compare(left.Interval, right.Interval)
		case "duration":
			if left.DurationMs < right.DurationMs {
				order = -1
			} else if left.DurationMs > right.DurationMs {
				order = 1
			} else {
				order = 0
			}
		default:
			order = compareTime(left.UpdatedAt, right.UpdatedAt)
		}
		if order == 0 {
			order = strings.Compare(strings.ToLower(left.Key), strings.ToLower(right.Key))
		}
		if ascending {
			return order < 0
		}
		return order > 0
	})
}

func buildMonitorKPI(monitors []MonitorSummary) MonitorKPI {
	kpi := MonitorKPI{
		Total: len(monitors),
	}
	for _, monitor := range monitors {
		switch monitor.Status {
		case "unhealthy":
			kpi.Unhealthy++
		case "disabled":
			kpi.Disabled++
		case "unknown":
			kpi.Unknown++
		}
	}
	return kpi
}

func collectMonitorGroups(monitors []MonitorSummary) []string {
	groupSet := make(map[string]struct{})
	for _, monitor := range monitors {
		group := strings.TrimSpace(monitor.Group)
		if group == "" {
			continue
		}
		groupSet[group] = struct{}{}
	}
	groups := make([]string, 0, len(groupSet))
	for group := range groupSet {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	return groups
}

func extractAlertTypes(alerts []*alert.Alert) []string {
	unique := make(map[string]struct{})
	for _, entry := range alerts {
		if entry == nil {
			continue
		}
		typeName := strings.TrimSpace(string(entry.Type))
		if len(typeName) == 0 {
			continue
		}
		unique[typeName] = struct{}{}
	}
	types := make([]string, 0, len(unique))
	for typeName := range unique {
		types = append(types, typeName)
	}
	sort.Strings(types)
	return types
}

func extractSuiteAlertTypes(endpoints []*endpoint.Endpoint) []string {
	unique := make(map[string]struct{})
	for _, endpointConfig := range endpoints {
		for _, entry := range endpointConfig.Alerts {
			if entry == nil {
				continue
			}
			typeName := strings.TrimSpace(string(entry.Type))
			if len(typeName) == 0 {
				continue
			}
			unique[typeName] = struct{}{}
		}
	}
	types := make([]string, 0, len(unique))
	for typeName := range unique {
		types = append(types, typeName)
	}
	sort.Strings(types)
	return types
}

func normalizeMonitorEntityType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case monitorEntityEndpoint, monitorEntitySuite, monitorEntityExternal:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return monitorEntityAll
	}
}

func normalizeMonitorEnabledFilter(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "false":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "all"
	}
}

func normalizeMonitorStatusFilter(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "healthy", "unhealthy", "unknown", "disabled":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "all"
	}
}

func normalizeMonitorGroupFilter(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || normalized == "all" {
		return ""
	}
	return normalized
}

func normalizeMonitorSortBy(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "name", "group", "status", "interval", "duration", "updatedat":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return "updatedat"
	}
}

func normalizeMonitorSortDirection(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "asc") {
		return "asc"
	}
	return "desc"
}

func parseMonitorPagination(pageValue, pageSizeValue string) (page, pageSize int) {
	page = 1
	pageSize = 50
	if parsed, err := strconv.Atoi(strings.TrimSpace(pageValue)); err == nil && parsed > 0 {
		page = parsed
	}
	if parsed, err := strconv.Atoi(strings.TrimSpace(pageSizeValue)); err == nil && parsed > 0 {
		pageSize = parsed
	}
	if pageSize > 200 {
		pageSize = 200
	}
	return page, pageSize
}
