package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/TwiN/gatus/v5/security"
)

func TestGroups(t *testing.T) {
	cfg := &config.Config{
		Endpoints:         []*endpoint.Endpoint{{Group: "core"}, {Group: "frontend"}},
		ExternalEndpoints: []*endpoint.ExternalEndpoint{{Group: "partner"}},
		Suites:            []*suite.Suite{{Group: "smoke"}, {Group: "regression"}},
		Security: &security.Config{Authorization: &security.AuthorizationConfig{
			EndpointGroups: []string{"core", "partner"},
			SuiteGroups:    []string{"smoke"},
		}},
	}
	api := New(cfg)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/groups", http.NoBody)
	response, err := api.Router().Test(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}
	var body GroupResponse
	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.EndpointGroups) != 2 || body.EndpointGroups[0] != "core" || body.EndpointGroups[1] != "partner" {
		t.Fatalf("unexpected endpoint groups: %+v", body.EndpointGroups)
	}
	if len(body.SuiteGroups) != 1 || body.SuiteGroups[0] != "smoke" {
		t.Fatalf("unexpected suite groups: %+v", body.SuiteGroups)
	}
}
