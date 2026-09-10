package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TwiN/gatus/v5/alerting"
	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/alerting/provider/custom"
	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/TwiN/gatus/v5/storage"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/watchdog"
)

func newSqliteStorageConfig(t *testing.T) *storage.Config {
	t.Helper()
	return &storage.Config{
		Type:                   storage.TypeSQLite,
		Path:                   filepath.Join(t.TempDir(), "gatus.db"),
		MaximumNumberOfResults: storage.DefaultMaximumNumberOfResults,
		MaximumNumberOfEvents:  storage.DefaultMaximumNumberOfEvents,
	}
}

func enabledBool() *bool {
	v := true
	return &v
}

func newSuiteWithAlert(name, group string, failureThreshold, successThreshold int) *suite.Suite {
	return &suite.Suite{
		Name:  name,
		Group: group,
		Endpoints: []*endpoint.Endpoint{
			{Name: "step-1", URL: "https://example.com", Interval: 1},
		},
		Alerts: []*alert.Alert{
			{
				Type:             alert.TypeCustom,
				Enabled:          enabledBool(),
				FailureThreshold: failureThreshold,
				SuccessThreshold: successThreshold,
				SendOnResolved:   enabledBool(),
				Description:      strPtr("suite-level"),
			},
		},
	}
}

func strPtr(s string) *string { return &s }

func TestInitializeStorage_ReloadsSuiteLevelTriggeredAlerts(t *testing.T) {
	storageCfg := newSqliteStorageConfig(t)
	s := newSuiteWithAlert("login-flow", "critical", 2, 2)
	s.Alerts[0].Triggered = true
	s.Alerts[0].ResolveKey = "rk-suite-1"
	s.NumberOfFailuresInARow = s.Alerts[0].FailureThreshold
	s.NumberOfSuccessesInARow = 1
	alertingEndpoint := s.ToEndpointForAlerting()

	if err := store.Initialize(storageCfg); err != nil {
		t.Fatalf("initialize store: %v", err)
	}
	if err := store.Get().UpsertTriggeredEndpointAlert(alertingEndpoint, s.Alerts[0]); err != nil {
		t.Fatalf("upsert triggered suite alert: %v", err)
	}
	store.Get().Close()

	// Simulate restart: fresh suite objects with Triggered=false and counters at 0.
	reloaded := newSuiteWithAlert("login-flow", "critical", 2, 2)
	cfg := &config.Config{
		Storage: storageCfg,
		Suites:  []*suite.Suite{reloaded},
	}
	initializeStorage(cfg)
	defer store.Get().Close()

	if !reloaded.Alerts[0].Triggered {
		t.Fatalf("expected suite alert Triggered=true after initializeStorage")
	}
	if reloaded.Alerts[0].ResolveKey != "rk-suite-1" {
		t.Fatalf("expected ResolveKey restored, got %q", reloaded.Alerts[0].ResolveKey)
	}
	if reloaded.NumberOfFailuresInARow != 2 {
		t.Fatalf("expected NumberOfFailuresInARow=2 (failure threshold), got %d", reloaded.NumberOfFailuresInARow)
	}
	if reloaded.NumberOfSuccessesInARow != 1 {
		t.Fatalf("expected NumberOfSuccessesInARow=1 restored, got %d", reloaded.NumberOfSuccessesInARow)
	}
}

func TestInitializeStorage_CleansOrphanSuiteTriggeredAlerts(t *testing.T) {
	storageCfg := newSqliteStorageConfig(t)
	orphan := newSuiteWithAlert("gone-suite", "critical", 2, 2)
	orphan.Alerts[0].Triggered = true
	orphan.Alerts[0].ResolveKey = "rk-orphan"
	alertingEndpoint := orphan.ToEndpointForAlerting()

	if err := store.Initialize(storageCfg); err != nil {
		t.Fatalf("initialize store: %v", err)
	}
	if err := store.Get().UpsertTriggeredEndpointAlert(alertingEndpoint, orphan.Alerts[0]); err != nil {
		t.Fatalf("upsert orphan suite alert: %v", err)
	}
	exists, _, _, err := store.Get().GetTriggeredEndpointAlert(alertingEndpoint, orphan.Alerts[0])
	if err != nil || !exists {
		t.Fatalf("expected orphan triggered alert to exist before re-init, exists=%v err=%v", exists, err)
	}
	store.Get().Close()

	// Config no longer contains the suite; synthetic suite key must be cleaned up.
	kept := newSuiteWithAlert("kept-suite", "critical", 2, 2)
	cfg := &config.Config{
		Storage: storageCfg,
		Suites:  []*suite.Suite{kept},
	}
	initializeStorage(cfg)
	defer store.Get().Close()

	exists, _, _, err = store.Get().GetTriggeredEndpointAlert(alertingEndpoint, orphan.Alerts[0])
	if err != nil {
		t.Fatalf("GetTriggeredEndpointAlert: %v", err)
	}
	if exists {
		t.Fatalf("expected orphan suite triggered alert to be deleted when suite key is gone")
	}
}

func TestInitializeStorage_CleansSuiteTriggeredAlertsWithChangedChecksum(t *testing.T) {
	storageCfg := newSqliteStorageConfig(t)
	s := newSuiteWithAlert("login-flow", "critical", 2, 2)
	s.Alerts[0].Triggered = true
	alertingEndpoint := s.ToEndpointForAlerting()

	if err := store.Initialize(storageCfg); err != nil {
		t.Fatalf("initialize store: %v", err)
	}
	if err := store.Get().UpsertTriggeredEndpointAlert(alertingEndpoint, s.Alerts[0]); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	store.Get().Close()

	// Same suite key, but alert configuration changed (different FailureThreshold → different checksum).
	changed := newSuiteWithAlert("login-flow", "critical", 5, 2)
	cfg := &config.Config{
		Storage: storageCfg,
		Suites:  []*suite.Suite{changed},
	}
	initializeStorage(cfg)
	defer store.Get().Close()

	if changed.Alerts[0].Triggered {
		t.Fatalf("expected Triggered=false when alert checksum no longer matches persisted row")
	}
	exists, _, _, err := store.Get().GetTriggeredEndpointAlert(changed.ToEndpointForAlerting(), changed.Alerts[0])
	if err != nil {
		t.Fatalf("GetTriggeredEndpointAlert: %v", err)
	}
	if exists {
		t.Fatalf("expected old checksum row cleaned; new checksum should not be triggered")
	}
	// Old alert checksum should also be gone.
	exists, _, _, err = store.Get().GetTriggeredEndpointAlert(alertingEndpoint, s.Alerts[0])
	if err != nil {
		t.Fatalf("GetTriggeredEndpointAlert(old): %v", err)
	}
	if exists {
		t.Fatalf("expected old checksum triggered alert deleted")
	}
}

func TestInitializeStorage_SuiteAlertResolveWorksAfterReload(t *testing.T) {
	_ = os.Setenv("MOCK_ALERT_PROVIDER", "true")
	t.Cleanup(func() { _ = os.Unsetenv("MOCK_ALERT_PROVIDER") })

	storageCfg := newSqliteStorageConfig(t)
	s := newSuiteWithAlert("login-flow", "critical", 2, 2)
	s.Alerts[0].Triggered = true
	s.Alerts[0].ResolveKey = "rk-resolve"
	s.NumberOfFailuresInARow = 2
	alertingEndpoint := s.ToEndpointForAlerting()

	if err := store.Initialize(storageCfg); err != nil {
		t.Fatalf("initialize store: %v", err)
	}
	if err := store.Get().UpsertTriggeredEndpointAlert(alertingEndpoint, s.Alerts[0]); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	store.Get().Close()

	reloaded := newSuiteWithAlert("login-flow", "critical", 2, 2)
	cfg := &config.Config{
		Storage: storageCfg,
		Suites:  []*suite.Suite{reloaded},
		Alerting: &alerting.Config{
			Custom: &custom.AlertProvider{
				DefaultConfig: custom.Config{URL: "https://twin.sh/health", Method: "GET"},
			},
		},
	}
	initializeStorage(cfg)
	defer store.Get().Close()

	if !reloaded.Alerts[0].Triggered {
		t.Fatalf("expected triggered after reload before resolve")
	}

	success := &suite.Result{
		Success: true,
		EndpointResults: []*endpoint.Result{
			{Name: "step-1", Success: true},
		},
	}
	watchdog.HandleSuiteAlerting(reloaded, success, cfg.Alerting)
	if reloaded.NumberOfSuccessesInARow != 1 || !reloaded.Alerts[0].Triggered {
		t.Fatalf("expected still triggered after 1 success, successes=%d triggered=%v", reloaded.NumberOfSuccessesInARow, reloaded.Alerts[0].Triggered)
	}
	watchdog.HandleSuiteAlerting(reloaded, success, cfg.Alerting)
	if reloaded.Alerts[0].Triggered {
		t.Fatalf("expected resolved after success threshold")
	}
	exists, _, _, err := store.Get().GetTriggeredEndpointAlert(reloaded.ToEndpointForAlerting(), reloaded.Alerts[0])
	if err != nil {
		t.Fatalf("GetTriggeredEndpointAlert: %v", err)
	}
	if exists {
		t.Fatalf("expected triggered alert row deleted after resolve")
	}
}
