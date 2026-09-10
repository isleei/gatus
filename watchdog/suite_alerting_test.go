package watchdog

import (
	"os"
	"testing"

	"github.com/TwiN/gatus/v5/alerting"
	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/alerting/provider/custom"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
)

func TestHandleSuiteAlerting(t *testing.T) {
	_ = os.Setenv("MOCK_ALERT_PROVIDER", "true")
	defer os.Clearenv()

	alertingConfig := &alerting.Config{
		Custom: &custom.AlertProvider{
			DefaultConfig: custom.Config{
				URL:    "https://twin.sh/health",
				Method: "GET",
			},
		},
	}
	enabled := true
	s := &suite.Suite{
		Name:  "login-flow",
		Group: "critical",
		Alerts: []*alert.Alert{
			{
				Type:             alert.TypeCustom,
				Enabled:          &enabled,
				FailureThreshold: 2,
				SuccessThreshold: 2,
				SendOnResolved:   &enabled,
			},
		},
	}

	fail := &suite.Result{
		Success: false,
		EndpointResults: []*endpoint.Result{
			{Name: "step-1", Success: false, Errors: []string{"status 500"}},
		},
		Errors: []string{"suite failed"},
	}
	success := &suite.Result{
		Success: true,
		EndpointResults: []*endpoint.Result{
			{Name: "step-1", Success: true},
		},
	}

	HandleSuiteAlerting(s, fail, alertingConfig)
	if s.NumberOfFailuresInARow != 1 || s.Alerts[0].Triggered {
		t.Fatalf("expected 1 failure and not triggered, got failures=%d triggered=%v", s.NumberOfFailuresInARow, s.Alerts[0].Triggered)
	}
	HandleSuiteAlerting(s, fail, alertingConfig)
	if s.NumberOfFailuresInARow != 2 || !s.Alerts[0].Triggered {
		t.Fatalf("expected alert triggered after 2 failures, got failures=%d triggered=%v", s.NumberOfFailuresInARow, s.Alerts[0].Triggered)
	}
	HandleSuiteAlerting(s, success, alertingConfig)
	if s.NumberOfSuccessesInARow != 1 || !s.Alerts[0].Triggered {
		t.Fatalf("expected still triggered after 1 success, got successes=%d triggered=%v", s.NumberOfSuccessesInARow, s.Alerts[0].Triggered)
	}
	HandleSuiteAlerting(s, success, alertingConfig)
	if s.NumberOfSuccessesInARow != 2 || s.Alerts[0].Triggered {
		t.Fatalf("expected resolved after 2 successes, got successes=%d triggered=%v", s.NumberOfSuccessesInARow, s.Alerts[0].Triggered)
	}
}

func TestHandleSuiteAlertingNilSafe(t *testing.T) {
	HandleSuiteAlerting(nil, nil, nil)
	HandleSuiteAlerting(&suite.Suite{Name: "x"}, &suite.Result{Success: true}, nil)
	HandleSuiteAlerting(&suite.Suite{Name: "x"}, &suite.Result{Success: true}, &alerting.Config{})
}
