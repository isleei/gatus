package security

import (
	"testing"

	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
)

func TestConfig_IsGroupAllowed(t *testing.T) {
	cfg := &Config{Authorization: &AuthorizationConfig{EndpointGroups: []string{"core", "backend"}, SuiteGroups: []string{"smoke"}}}
	if !cfg.IsEndpointGroupAllowed("CORE") {
		t.Fatal("expected endpoint group to be allowed")
	}
	if cfg.IsEndpointGroupAllowed("frontend") {
		t.Fatal("expected endpoint group to be denied")
	}
	if !cfg.IsSuiteGroupAllowed("smoke") {
		t.Fatal("expected suite group to be allowed")
	}
	if cfg.IsSuiteGroupAllowed("regression") {
		t.Fatal("expected suite group to be denied")
	}
}

func TestConfig_FilterStatuses(t *testing.T) {
	cfg := &Config{Authorization: &AuthorizationConfig{EndpointGroups: []string{"core"}, SuiteGroups: []string{"smoke"}}}
	endpointStatuses := []*endpoint.Status{{Group: "core"}, {Group: "frontend"}}
	suiteStatuses := []*suite.Status{{Group: "smoke"}, {Group: "regression"}}

	filteredEndpointStatuses := cfg.FilterEndpointStatuses(endpointStatuses)
	if len(filteredEndpointStatuses) != 1 || filteredEndpointStatuses[0].Group != "core" {
		t.Fatalf("unexpected endpoint statuses: %+v", filteredEndpointStatuses)
	}

	filteredSuiteStatuses := cfg.FilterSuiteStatuses(suiteStatuses)
	if len(filteredSuiteStatuses) != 1 || filteredSuiteStatuses[0].Group != "smoke" {
		t.Fatalf("unexpected suite statuses: %+v", filteredSuiteStatuses)
	}
}
