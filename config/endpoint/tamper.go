package endpoint

import (
	"errors"
	"strings"
)

const (
	DefaultTamperBaselineSamples     = 20
	DefaultTamperDriftThresholdPct   = 20
	DefaultTamperConsecutiveBreaches = 3
)

var (
	ErrInvalidTamperBaselineSamples   = errors.New("tamper baseline-samples must be greater than 0")
	ErrInvalidTamperDriftThresholdPct = errors.New("tamper drift-threshold-percent must be greater than 0")
	ErrInvalidTamperConsecutive       = errors.New("tamper consecutive-breaches must be greater than 0")
)

// TamperConfig defines body-size drift based tamper detection parameters.
type TamperConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`

	// BaselineSamples is the number of recent samples used to compute the body size baseline.
	BaselineSamples int `yaml:"baseline-samples,omitempty"`

	// DriftThresholdPercent is the allowed drift percentage before counting a breach.
	DriftThresholdPercent int64 `yaml:"drift-threshold-percent,omitempty"`

	// ConsecutiveBreaches is the number of consecutive breaches required to fail a check.
	ConsecutiveBreaches int `yaml:"consecutive-breaches,omitempty"`

	// RequiredSubstrings is a list of strings that must be present in the response body.
	RequiredSubstrings []string `yaml:"required-substrings,omitempty"`

	// ForbiddenSubstrings is a list of strings that must not be present in the response body.
	ForbiddenSubstrings []string `yaml:"forbidden-substrings,omitempty"`
}

// ValidateAndSetDefaults validates tamper configuration and applies defaults when enabled.
func (config *TamperConfig) ValidateAndSetDefaults() error {
	if config == nil || !config.Enabled {
		return nil
	}
	if config.BaselineSamples == 0 {
		config.BaselineSamples = DefaultTamperBaselineSamples
	}
	if config.DriftThresholdPercent == 0 {
		config.DriftThresholdPercent = DefaultTamperDriftThresholdPct
	}
	if config.ConsecutiveBreaches == 0 {
		config.ConsecutiveBreaches = DefaultTamperConsecutiveBreaches
	}
	if config.BaselineSamples < 1 {
		return ErrInvalidTamperBaselineSamples
	}
	if config.DriftThresholdPercent < 1 {
		return ErrInvalidTamperDriftThresholdPct
	}
	if config.ConsecutiveBreaches < 1 {
		return ErrInvalidTamperConsecutive
	}
	config.RequiredSubstrings = normalizeTamperSubstrings(config.RequiredSubstrings)
	config.ForbiddenSubstrings = normalizeTamperSubstrings(config.ForbiddenSubstrings)
	return nil
}

func normalizeTamperSubstrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if len(trimmed) == 0 {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	if len(normalized) == 0 {
		return nil
	}
	return normalized
}
