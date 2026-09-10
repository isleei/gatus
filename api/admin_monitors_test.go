package api

import (
	"reflect"
	"testing"
	"time"
)

func sampleMonitors() []MonitorSummary {
	t1 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC)
	return []MonitorSummary{
		{
			Key:               "core_frontend",
			EntityType:        monitorEntityEndpoint,
			Type:              "HTTP",
			Name:              "Frontend",
			Group:             "Core",
			Enabled:           true,
			Status:            "healthy",
			URL:               "https://frontend.example",
			Interval:          "1m0s",
			DurationMs:        120,
			NotificationTypes: []string{"slack"},
			UpdatedAt:         &t2,
		},
		{
			Key:        "core_api",
			EntityType: monitorEntityEndpoint,
			Type:       "HTTP",
			Name:       "API",
			Group:      "Core",
			Enabled:    true,
			Status:     "unhealthy",
			URL:        "https://api.example",
			Interval:   "30s",
			DurationMs: 500,
			UpdatedAt:  &t3,
		},
		{
			Key:        "ops_backup",
			EntityType: monitorEntitySuite,
			Type:       "SUITE",
			Name:       "Backup Suite",
			Group:      "Ops",
			Enabled:    false,
			Status:     "disabled",
			Interval:   "5m0s",
			DurationMs: 0,
		},
		{
			Key:        "ext_heartbeat",
			EntityType: monitorEntityExternal,
			Type:       "EXTERNAL",
			Name:       "Heartbeat",
			Group:      "",
			Enabled:    true,
			Status:     "unknown",
			URL:        "external://ext_heartbeat",
			Interval:   "1m0s",
			DurationMs: 10,
			UpdatedAt:  &t1,
		},
	}
}

func monitorKeys(monitors []MonitorSummary) []string {
	if monitors == nil {
		return nil
	}
	keys := make([]string, len(monitors))
	for i, monitor := range monitors {
		keys[i] = monitor.Key
	}
	return keys
}

func TestFilterMonitors(t *testing.T) {
	monitors := sampleMonitors()

	testCases := []struct {
		name          string
		input         []MonitorSummary
		entityType    string
		query         string
		groupFilter   string
		enabledFilter string
		statusFilter  string
		wantKeys      []string
	}{
		{
			name:          "nil input returns empty non-nil slice",
			input:         nil,
			entityType:    monitorEntityAll,
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{},
		},
		{
			name:          "empty input returns empty slice",
			input:         []MonitorSummary{},
			entityType:    monitorEntityAll,
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{},
		},
		{
			name:          "entityType endpoint",
			input:         monitors,
			entityType:    monitorEntityEndpoint,
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"core_frontend", "core_api"},
		},
		{
			name:          "entityType suite",
			input:         monitors,
			entityType:    monitorEntitySuite,
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"ops_backup"},
		},
		{
			name:          "entityType external",
			input:         monitors,
			entityType:    monitorEntityExternal,
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"ext_heartbeat"},
		},
		{
			name:          "group filter case-insensitive",
			input:         monitors,
			entityType:    monitorEntityAll,
			groupFilter:   "core",
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"core_frontend", "core_api"},
		},
		{
			name:          "group filter all is no-op",
			input:         monitors,
			entityType:    monitorEntityAll,
			groupFilter:   "all",
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"core_frontend", "core_api", "ops_backup", "ext_heartbeat"},
		},
		{
			name:          "enabled true",
			input:         monitors,
			entityType:    monitorEntityAll,
			enabledFilter: "true",
			statusFilter:  "all",
			wantKeys:      []string{"core_frontend", "core_api", "ext_heartbeat"},
		},
		{
			name:          "enabled false",
			input:         monitors,
			entityType:    monitorEntityAll,
			enabledFilter: "false",
			statusFilter:  "all",
			wantKeys:      []string{"ops_backup"},
		},
		{
			name:          "status unhealthy",
			input:         monitors,
			entityType:    monitorEntityAll,
			enabledFilter: "all",
			statusFilter:  "unhealthy",
			wantKeys:      []string{"core_api"},
		},
		{
			name:          "status unknown",
			input:         monitors,
			entityType:    monitorEntityAll,
			enabledFilter: "all",
			statusFilter:  "unknown",
			wantKeys:      []string{"ext_heartbeat"},
		},
		{
			name:          "q matches name",
			input:         monitors,
			entityType:    monitorEntityAll,
			query:         "backup",
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"ops_backup"},
		},
		{
			name:          "q matches url",
			input:         monitors,
			entityType:    monitorEntityAll,
			query:         "api.example",
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"core_api"},
		},
		{
			name:          "q matches notification type",
			input:         monitors,
			entityType:    monitorEntityAll,
			query:         "slack",
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{"core_frontend"},
		},
		{
			name:          "combined entity group and status",
			input:         monitors,
			entityType:    monitorEntityEndpoint,
			groupFilter:   "Core",
			enabledFilter: "true",
			statusFilter:  "healthy",
			wantKeys:      []string{"core_frontend"},
		},
		{
			name:          "q no match",
			input:         monitors,
			entityType:    monitorEntityAll,
			query:         "does-not-exist",
			enabledFilter: "all",
			statusFilter:  "all",
			wantKeys:      []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := filterMonitors(testCase.input, testCase.entityType, testCase.query, testCase.groupFilter, testCase.enabledFilter, testCase.statusFilter)
			if got == nil {
				t.Fatal("expected non-nil filtered slice")
			}
			if !reflect.DeepEqual(monitorKeys(got), testCase.wantKeys) {
				t.Fatalf("expected keys %v, got %v", testCase.wantKeys, monitorKeys(got))
			}
		})
	}
}

func TestSortMonitors(t *testing.T) {
	base := sampleMonitors()

	testCases := []struct {
		name     string
		sortBy   string
		sortDir  string
		wantKeys []string
	}{
		{
			name:     "name asc",
			sortBy:   "name",
			sortDir:  "asc",
			wantKeys: []string{"core_api", "ops_backup", "core_frontend", "ext_heartbeat"},
		},
		{
			name:     "name desc",
			sortBy:   "name",
			sortDir:  "desc",
			wantKeys: []string{"ext_heartbeat", "core_frontend", "ops_backup", "core_api"},
		},
		{
			name:     "group asc empty first then key tie-break",
			sortBy:   "group",
			sortDir:  "asc",
			wantKeys: []string{"ext_heartbeat", "core_api", "core_frontend", "ops_backup"},
		},
		{
			name:     "status asc by severity weight",
			sortBy:   "status",
			sortDir:  "asc",
			wantKeys: []string{"core_api", "ext_heartbeat", "core_frontend", "ops_backup"},
		},
		{
			name:     "interval asc",
			sortBy:   "interval",
			sortDir:  "asc",
			wantKeys: []string{"core_frontend", "ext_heartbeat", "core_api", "ops_backup"},
		},
		{
			name:     "duration desc",
			sortBy:   "duration",
			sortDir:  "desc",
			wantKeys: []string{"core_api", "core_frontend", "ext_heartbeat", "ops_backup"},
		},
		{
			name:     "updatedAt desc nil last when descending",
			sortBy:   "updatedat",
			sortDir:  "desc",
			wantKeys: []string{"core_api", "core_frontend", "ext_heartbeat", "ops_backup"},
		},
		{
			name:     "updatedAt asc nil first",
			sortBy:   "updatedat",
			sortDir:  "asc",
			wantKeys: []string{"ops_backup", "ext_heartbeat", "core_frontend", "core_api"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			monitors := append([]MonitorSummary(nil), base...)
			sortMonitors(monitors, testCase.sortBy, testCase.sortDir)
			if !reflect.DeepEqual(monitorKeys(monitors), testCase.wantKeys) {
				t.Fatalf("expected keys %v, got %v", testCase.wantKeys, monitorKeys(monitors))
			}
		})
	}

	t.Run("nil and empty slices are safe", func(t *testing.T) {
		sortMonitors(nil, "name", "asc")
		empty := []MonitorSummary{}
		sortMonitors(empty, "name", "asc")
		if len(empty) != 0 {
			t.Fatalf("expected empty slice to remain empty, got %d", len(empty))
		}
	})
}

func TestBuildMonitorKPI(t *testing.T) {
	testCases := []struct {
		name     string
		input    []MonitorSummary
		expected MonitorKPI
	}{
		{
			name:     "nil slice",
			input:    nil,
			expected: MonitorKPI{Total: 0},
		},
		{
			name:     "empty slice",
			input:    []MonitorSummary{},
			expected: MonitorKPI{Total: 0},
		},
		{
			name:  "counts unhealthy disabled unknown only",
			input: sampleMonitors(),
			expected: MonitorKPI{
				Total:     4,
				Unhealthy: 1,
				Disabled:  1,
				Unknown:   1,
			},
		},
		{
			name: "healthy alone does not inflate KPI buckets",
			input: []MonitorSummary{
				{Key: "a", Status: "healthy"},
				{Key: "b", Status: "healthy"},
			},
			expected: MonitorKPI{Total: 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := buildMonitorKPI(testCase.input)
			if !reflect.DeepEqual(got, testCase.expected) {
				t.Fatalf("expected %+v, got %+v", testCase.expected, got)
			}
		})
	}
}

func TestParseMonitorPagination(t *testing.T) {
	testCases := []struct {
		name         string
		page         string
		pageSize     string
		wantPage     int
		wantPageSize int
	}{
		{name: "defaults", page: "", pageSize: "", wantPage: 1, wantPageSize: 50},
		{name: "whitespace defaults", page: "  ", pageSize: "\t", wantPage: 1, wantPageSize: 50},
		{name: "valid values", page: "3", pageSize: "25", wantPage: 3, wantPageSize: 25},
		{name: "zero falls back", page: "0", pageSize: "0", wantPage: 1, wantPageSize: 50},
		{name: "negative falls back", page: "-2", pageSize: "-10", wantPage: 1, wantPageSize: 50},
		{name: "non-numeric falls back", page: "abc", pageSize: "xyz", wantPage: 1, wantPageSize: 50},
		{name: "pageSize capped at 200", page: "1", pageSize: "500", wantPage: 1, wantPageSize: 200},
		{name: "pageSize exactly 200", page: "2", pageSize: "200", wantPage: 2, wantPageSize: 200},
		{name: "trimmed numeric", page: " 4 ", pageSize: " 10 ", wantPage: 4, wantPageSize: 10},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			page, pageSize := parseMonitorPagination(testCase.page, testCase.pageSize)
			if page != testCase.wantPage || pageSize != testCase.wantPageSize {
				t.Fatalf("expected (%d, %d), got (%d, %d)", testCase.wantPage, testCase.wantPageSize, page, pageSize)
			}
		})
	}
}
