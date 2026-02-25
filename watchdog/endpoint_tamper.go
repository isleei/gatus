package watchdog

import (
	"fmt"
	"math"
	"strings"

	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common/paging"
	"github.com/TwiN/logr"
)

func applyBodySizeTamperDetection(ep *endpoint.Endpoint, result *endpoint.Result) {
	if ep == nil || result == nil || ep.TamperConfig == nil || !ep.TamperConfig.Enabled {
		return
	}
	bodySize := int64(len(result.Body))
	result.BodySizeBytes = &bodySize
	applyBodyContentTamperChecks(ep, result)

	// Skip active tamper evaluation on connectivity/execution errors.
	if len(result.Errors) > 0 || !result.Connected {
		ep.NumberOfBodySizeDriftBreachesInARow = 0
		return
	}

	numberOfSamples := ep.TamperConfig.BaselineSamples
	if numberOfSamples < 1 {
		numberOfSamples = endpoint.DefaultTamperBaselineSamples
	}
	driftThresholdPercent := ep.TamperConfig.DriftThresholdPercent
	if driftThresholdPercent < 1 {
		driftThresholdPercent = endpoint.DefaultTamperDriftThresholdPct
	}
	consecutiveBreaches := ep.TamperConfig.ConsecutiveBreaches
	if consecutiveBreaches < 1 {
		consecutiveBreaches = endpoint.DefaultTamperConsecutiveBreaches
	}
	status, err := store.Get().GetEndpointStatusByKey(ep.Key(), paging.NewEndpointStatusParams().WithResults(1, numberOfSamples))
	if err != nil {
		logr.Errorf("[watchdog.applyBodySizeTamperDetection] Failed loading history for key=%s: %s", ep.Key(), err.Error())
		ep.NumberOfBodySizeDriftBreachesInARow = 0
		return
	}
	if status == nil || len(status.Results) == 0 {
		// Warm-up: no history yet, only persist observations.
		ep.NumberOfBodySizeDriftBreachesInARow = 0
		return
	}

	var bodySizeSum int64
	var usableSamples int
	for _, historicalResult := range status.Results {
		if historicalResult == nil || historicalResult.BodySizeBytes == nil {
			continue
		}
		bodySizeSum += *historicalResult.BodySizeBytes
		usableSamples++
	}
	if usableSamples == 0 {
		// Warm-up: no historical row has body size observations yet.
		ep.NumberOfBodySizeDriftBreachesInARow = 0
		return
	}

	bodySizeBaseline := int64(math.Round(float64(bodySizeSum) / float64(usableSamples)))
	result.BodySizeBaselineBytes = &bodySizeBaseline
	driftPercent := calculateBodySizeDriftPercent(bodySize, bodySizeBaseline)
	result.BodySizeDriftPercent = &driftPercent

	// Warm-up: baseline is displayed, but failures are not triggered until enough samples are available.
	if usableSamples < numberOfSamples {
		ep.NumberOfBodySizeDriftBreachesInARow = 0
		return
	}

	if driftPercent <= float64(driftThresholdPercent) {
		ep.NumberOfBodySizeDriftBreachesInARow = 0
		return
	}
	ep.NumberOfBodySizeDriftBreachesInARow++
	if ep.NumberOfBodySizeDriftBreachesInARow < consecutiveBreaches {
		return
	}

	result.Success = false
	result.ConditionResults = append(result.ConditionResults, &endpoint.ConditionResult{
		Condition: fmt.Sprintf(
			"[BODY_SIZE_DRIFT] <= %d%% (actual: %.0f%%)",
			driftThresholdPercent,
			driftPercent,
		),
		Success: false,
	})
}

func applyBodyContentTamperChecks(ep *endpoint.Endpoint, result *endpoint.Result) {
	if ep == nil || result == nil || ep.TamperConfig == nil || !ep.TamperConfig.Enabled {
		return
	}
	if len(ep.TamperConfig.RequiredSubstrings) == 0 && len(ep.TamperConfig.ForbiddenSubstrings) == 0 {
		return
	}
	body := strings.ToLower(string(result.Body))
	for _, requiredSubstring := range ep.TamperConfig.RequiredSubstrings {
		if strings.Contains(body, strings.ToLower(requiredSubstring)) {
			continue
		}
		result.Success = false
		result.ConditionResults = append(result.ConditionResults, &endpoint.ConditionResult{
			Condition: fmt.Sprintf(`[BODY_REQUIRED_SUBSTRING] contains "%s"`, requiredSubstring),
			Success:   false,
		})
	}
	for _, forbiddenSubstring := range ep.TamperConfig.ForbiddenSubstrings {
		if !strings.Contains(body, strings.ToLower(forbiddenSubstring)) {
			continue
		}
		result.Success = false
		result.ConditionResults = append(result.ConditionResults, &endpoint.ConditionResult{
			Condition: fmt.Sprintf(`[BODY_FORBIDDEN_SUBSTRING] not contains "%s"`, forbiddenSubstring),
			Success:   false,
		})
	}
}

func calculateBodySizeDriftPercent(currentBodySize, baselineBodySize int64) float64 {
	if baselineBodySize == 0 {
		if currentBodySize == 0 {
			return 0
		}
		return 100
	}
	return math.Abs(float64(currentBodySize-baselineBodySize)) * 100 / float64(baselineBodySize)
}
