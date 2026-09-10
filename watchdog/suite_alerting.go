package watchdog

import (
	"fmt"
	"strings"

	"github.com/TwiN/gatus/v5/alerting"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
)

// HandleSuiteAlerting triggers or resolves suite-level alerts based on suite execution result.
// Providers receive a synthetic endpoint derived from the suite so existing WeCom/Slack/custom paths work.
func HandleSuiteAlerting(s *suite.Suite, result *suite.Result, alertingConfig *alerting.Config) {
	if s == nil || len(s.Alerts) == 0 || alertingConfig == nil {
		return
	}
	ep := s.ToEndpointForAlerting()
	epResult := suiteResultToEndpointResult(s, result)
	HandleAlerting(ep, epResult, alertingConfig)
	s.NumberOfFailuresInARow = ep.NumberOfFailuresInARow
	s.NumberOfSuccessesInARow = ep.NumberOfSuccessesInARow
	s.LastReminderSent = ep.LastReminderSent
}

func suiteResultToEndpointResult(s *suite.Suite, result *suite.Result) *endpoint.Result {
	if result == nil {
		return &endpoint.Result{Success: false, Errors: []string{"nil suite result"}}
	}
	epResult := &endpoint.Result{
		Name:      s.Name,
		Success:   result.Success,
		Timestamp: result.Timestamp,
		Duration:  result.Duration,
		Errors:    append([]string{}, result.Errors...),
	}
	var failedSteps []string
	for _, step := range result.EndpointResults {
		if step == nil {
			continue
		}
		if step.Success {
			epResult.ConditionResults = append(epResult.ConditionResults, &endpoint.ConditionResult{
				Condition: fmt.Sprintf("suite step %s", step.Name),
				Success:   true,
			})
			continue
		}
		failedSteps = append(failedSteps, step.Name)
		epResult.ConditionResults = append(epResult.ConditionResults, &endpoint.ConditionResult{
			Condition: fmt.Sprintf("suite step %s", step.Name),
			Success:   false,
		})
		for _, errMsg := range step.Errors {
			epResult.Errors = append(epResult.Errors, fmt.Sprintf("%s: %s", step.Name, errMsg))
		}
	}
	if len(failedSteps) > 0 {
		epResult.Errors = append(epResult.Errors, "failed steps: "+strings.Join(failedSteps, ", "))
	}
	return epResult
}
