package endpoint

import "testing"

func TestTamperConfigValidateAndSetDefaults(t *testing.T) {
	config := &TamperConfig{
		Enabled:             true,
		RequiredSubstrings:  []string{"  <title>Status</title>  ", "", "在线状态页"},
		ForbiddenSubstrings: []string{"  博彩  ", "   "},
	}
	if err := config.ValidateAndSetDefaults(); err != nil {
		t.Fatalf("did not expect error, got %v", err)
	}
	if config.BaselineSamples != DefaultTamperBaselineSamples {
		t.Fatalf("expected baseline samples default=%d, got %d", DefaultTamperBaselineSamples, config.BaselineSamples)
	}
	if config.DriftThresholdPercent != DefaultTamperDriftThresholdPct {
		t.Fatalf("expected drift threshold default=%d, got %d", DefaultTamperDriftThresholdPct, config.DriftThresholdPercent)
	}
	if config.ConsecutiveBreaches != DefaultTamperConsecutiveBreaches {
		t.Fatalf("expected consecutive breaches default=%d, got %d", DefaultTamperConsecutiveBreaches, config.ConsecutiveBreaches)
	}
	if len(config.RequiredSubstrings) != 2 {
		t.Fatalf("expected 2 normalized required substrings, got %v", config.RequiredSubstrings)
	}
	if len(config.ForbiddenSubstrings) != 1 {
		t.Fatalf("expected 1 normalized forbidden substring, got %v", config.ForbiddenSubstrings)
	}
}

func TestTamperConfigValidateAndSetDefaultsInvalid(t *testing.T) {
	scenarios := []struct {
		name      string
		config    *TamperConfig
		validator error
	}{
		{
			name: "invalid-baseline-samples",
			config: &TamperConfig{
				Enabled:         true,
				BaselineSamples: -1,
			},
			validator: ErrInvalidTamperBaselineSamples,
		},
		{
			name: "invalid-drift-threshold-percent",
			config: &TamperConfig{
				Enabled:               true,
				BaselineSamples:       1,
				DriftThresholdPercent: -1,
			},
			validator: ErrInvalidTamperDriftThresholdPct,
		},
		{
			name: "invalid-consecutive-breaches",
			config: &TamperConfig{
				Enabled:               true,
				BaselineSamples:       1,
				DriftThresholdPercent: 1,
				ConsecutiveBreaches:   -1,
			},
			validator: ErrInvalidTamperConsecutive,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			if err := scenario.config.ValidateAndSetDefaults(); err != scenario.validator {
				t.Fatalf("expected error %v, got %v", scenario.validator, err)
			}
		})
	}
}
