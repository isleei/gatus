package api

import (
	"testing"

	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
)

func TestApplyBatchOperationToSuiteSetAlertTypesCreatesIndependentAlertsPerEndpoint(t *testing.T) {
	monitoredSuite := &suite.Suite{
		Endpoints: []*endpoint.Endpoint{
			{Name: "step-1", URL: "https://example.org/step-1"},
			{Name: "step-2", URL: "https://example.org/step-2"},
		},
	}

	err := applyBatchOperationToSuite(monitoredSuite, "set-alert-types", map[string]interface{}{
		"alertTypes": []interface{}{"slack"},
	})
	if err != nil {
		t.Fatalf("expected no error applying batch alert type update, got: %v", err)
	}

	if len(monitoredSuite.Endpoints[0].Alerts) != 1 || len(monitoredSuite.Endpoints[1].Alerts) != 1 {
		t.Fatalf("expected one alert on each suite endpoint, got %d and %d", len(monitoredSuite.Endpoints[0].Alerts), len(monitoredSuite.Endpoints[1].Alerts))
	}

	if monitoredSuite.Endpoints[0].Alerts[0] == monitoredSuite.Endpoints[1].Alerts[0] {
		t.Fatal("expected suite endpoints to hold independent alert instances")
	}

	monitoredSuite.Endpoints[0].Alerts[0].Triggered = true
	if monitoredSuite.Endpoints[1].Alerts[0].Triggered {
		t.Fatal("expected mutating one endpoint alert state not to affect another endpoint")
	}
}
