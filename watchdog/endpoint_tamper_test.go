package watchdog

import (
	"strings"
	"testing"
	"time"

	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/storage/store"
)

func TestCalculateBodySizeDriftPercent(t *testing.T) {
	if value := calculateBodySizeDriftPercent(0, 0); value != 0 {
		t.Fatalf("expected 0 drift for zero baseline and zero body, got %f", value)
	}
	if value := calculateBodySizeDriftPercent(128, 0); value != 100 {
		t.Fatalf("expected 100 drift for non-zero body and zero baseline, got %f", value)
	}
	if value := calculateBodySizeDriftPercent(120, 100); value != 20 {
		t.Fatalf("expected 20 drift, got %f", value)
	}
}

func TestApplyBodySizeTamperDetectionWarmup(t *testing.T) {
	store.Get().Clear()
	ep := &endpoint.Endpoint{
		Name:       "warmup",
		Group:      "tamper",
		URL:        "https://example.org",
		Conditions: []endpoint.Condition{"[STATUS] == 200"},
		TamperConfig: &endpoint.TamperConfig{
			Enabled:               true,
			BaselineSamples:       3,
			DriftThresholdPercent: 20,
			ConsecutiveBreaches:   2,
		},
	}
	_ = store.Get().InsertEndpointResult(ep, &endpoint.Result{
		Success:       true,
		Connected:     true,
		BodySizeBytes: int64Ptr(100),
		Timestamp:     time.Now().Add(-3 * time.Minute),
	})
	_ = store.Get().InsertEndpointResult(ep, &endpoint.Result{
		Success:       true,
		Connected:     true,
		BodySizeBytes: int64Ptr(101),
		Timestamp:     time.Now().Add(-2 * time.Minute),
	})

	result := &endpoint.Result{Success: true, Connected: true, Body: make([]byte, 120)}
	applyBodySizeTamperDetection(ep, result)

	if !result.Success {
		t.Fatal("expected warm-up execution to remain healthy")
	}
	if result.BodySizeBytes == nil || *result.BodySizeBytes != 120 {
		t.Fatalf("expected body size bytes to be 120, got %+v", result.BodySizeBytes)
	}
	if result.BodySizeBaselineBytes == nil || *result.BodySizeBaselineBytes != 101 {
		t.Fatalf("expected warm-up baseline to be 101, got %+v", result.BodySizeBaselineBytes)
	}
	if result.BodySizeDriftPercent == nil || *result.BodySizeDriftPercent < 18 {
		t.Fatalf("expected warm-up drift around 18%%, got %+v", result.BodySizeDriftPercent)
	}
	if ep.NumberOfBodySizeDriftBreachesInARow != 0 {
		t.Fatalf("expected consecutive breach counter reset to 0, got %d", ep.NumberOfBodySizeDriftBreachesInARow)
	}
}

func TestApplyBodySizeTamperDetectionConsecutiveBreaches(t *testing.T) {
	store.Get().Clear()
	ep := &endpoint.Endpoint{
		Name:       "breaches",
		Group:      "tamper",
		URL:        "https://example.org",
		Conditions: []endpoint.Condition{"[STATUS] == 200"},
		TamperConfig: &endpoint.TamperConfig{
			Enabled:               true,
			BaselineSamples:       3,
			DriftThresholdPercent: 20,
			ConsecutiveBreaches:   2,
		},
	}
	_ = store.Get().InsertEndpointResult(ep, &endpoint.Result{
		Success:       true,
		Connected:     true,
		BodySizeBytes: int64Ptr(100),
		Timestamp:     time.Now().Add(-5 * time.Minute),
	})
	_ = store.Get().InsertEndpointResult(ep, &endpoint.Result{
		Success:       true,
		Connected:     true,
		BodySizeBytes: int64Ptr(99),
		Timestamp:     time.Now().Add(-4 * time.Minute),
	})
	_ = store.Get().InsertEndpointResult(ep, &endpoint.Result{
		Success:       true,
		Connected:     true,
		BodySizeBytes: int64Ptr(101),
		Timestamp:     time.Now().Add(-3 * time.Minute),
	})

	first := &endpoint.Result{Success: true, Connected: true, Body: make([]byte, 140)}
	applyBodySizeTamperDetection(ep, first)
	if !first.Success {
		t.Fatal("expected first breach not to fail due to consecutive threshold")
	}
	if ep.NumberOfBodySizeDriftBreachesInARow != 1 {
		t.Fatalf("expected counter=1 after first breach, got %d", ep.NumberOfBodySizeDriftBreachesInARow)
	}

	second := &endpoint.Result{Success: true, Connected: true, Body: make([]byte, 150)}
	applyBodySizeTamperDetection(ep, second)
	if second.Success {
		t.Fatal("expected second consecutive breach to fail")
	}
	if second.BodySizeBaselineBytes == nil || *second.BodySizeBaselineBytes != 100 {
		t.Fatalf("expected baseline=100, got %+v", second.BodySizeBaselineBytes)
	}
	if second.BodySizeDriftPercent == nil || *second.BodySizeDriftPercent < 49 {
		t.Fatalf("expected drift around 50%%, got %+v", second.BodySizeDriftPercent)
	}
	if len(second.ConditionResults) != 1 || !strings.Contains(second.ConditionResults[0].Condition, "[BODY_SIZE_DRIFT]") {
		t.Fatalf("expected body size drift condition result, got %+v", second.ConditionResults)
	}
}

func TestApplyBodySizeTamperDetectionWithBodyContentRules(t *testing.T) {
	store.Get().Clear()
	ep := &endpoint.Endpoint{
		Name:       "content-rules",
		Group:      "tamper",
		URL:        "https://example.org",
		Conditions: []endpoint.Condition{"[STATUS] == 200"},
		TamperConfig: &endpoint.TamperConfig{
			Enabled:               true,
			BaselineSamples:       1,
			DriftThresholdPercent: 20,
			ConsecutiveBreaches:   1,
			RequiredSubstrings:    []string{"在线状态页", "<title>状态页</title>"},
			ForbiddenSubstrings:   []string{"博彩", "挂马"},
		},
	}
	_ = store.Get().InsertEndpointResult(ep, &endpoint.Result{
		Success:       true,
		Connected:     true,
		BodySizeBytes: int64Ptr(80),
		Timestamp:     time.Now().Add(-time.Minute),
	})

	result := &endpoint.Result{
		Success:   true,
		Connected: true,
		Body:      []byte("<html><title>被篡改</title><body>博彩广告</body></html>"),
	}
	applyBodySizeTamperDetection(ep, result)
	if result.Success {
		t.Fatal("expected content rule failures to mark result unhealthy")
	}
	if len(result.ConditionResults) < 2 {
		t.Fatalf("expected at least 2 content tamper condition results, got %+v", result.ConditionResults)
	}
	hasRequiredFailure := false
	hasForbiddenFailure := false
	for _, conditionResult := range result.ConditionResults {
		if strings.Contains(conditionResult.Condition, "[BODY_REQUIRED_SUBSTRING]") {
			hasRequiredFailure = true
		}
		if strings.Contains(conditionResult.Condition, "[BODY_FORBIDDEN_SUBSTRING]") {
			hasForbiddenFailure = true
		}
	}
	if !hasRequiredFailure {
		t.Fatalf("expected required substring failure in condition results, got %+v", result.ConditionResults)
	}
	if !hasForbiddenFailure {
		t.Fatalf("expected forbidden substring failure in condition results, got %+v", result.ConditionResults)
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}
