package ui

import (
	"errors"
	"testing"
)

func TestValidateAndSetDefaults(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr error
	}{
		{
			name: "with-valid-config",
			config: &Config{
				Badge: &Badge{
					ResponseTime: &ResponseTime{Thresholds: []int{50, 200, 300, 500, 750}},
				},
			},
			wantErr: nil,
		},
		{
			name: "with-invalid-threshold-length",
			config: &Config{
				Badge: &Badge{
					ResponseTime: &ResponseTime{Thresholds: []int{50, 200, 300, 500}},
				},
			},
			wantErr: ErrInvalidBadgeResponseTimeConfig,
		},
		{
			name: "with-invalid-thresholds-order",
			config: &Config{
				Badge: &Badge{ResponseTime: &ResponseTime{Thresholds: []int{50, 200, 500, 300, 750}}},
			},
			wantErr: ErrInvalidBadgeResponseTimeConfig,
		},
		{
			name:    "with-no-badge-configured", // should give default badge cfg
			config:  &Config{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.ValidateAndSetDefaults(); !errors.Is(err, tt.wantErr) {
				t.Errorf("Expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestConfig_Clone(t *testing.T) {
	original := &Config{
		HideErrors: true,
		Badge: &Badge{
			ResponseTime: &ResponseTime{Thresholds: []int{10, 20, 30, 40, 50}},
		},
	}

	clone := original.Clone()
	if clone == nil {
		t.Fatal("expected clone to be non-nil")
	}
	if !clone.HideErrors {
		t.Fatal("expected clone to copy HideErrors")
	}
	clone.Badge.ResponseTime.Thresholds[0] = 999
	if original.Badge.ResponseTime.Thresholds[0] == 999 {
		t.Fatal("expected clone thresholds to be deep copied")
	}
}
