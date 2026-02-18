package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/alerting"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
)

const (
	// GatusManagedOverlayPathEnvVar is the environment variable used to override the managed overlay path.
	GatusManagedOverlayPathEnvVar = "GATUS_MANAGED_OVERLAY_PATH"
)

// ManagedOverlay contains configuration sections that can be fully managed by the web UI.
type ManagedOverlay struct {
	Alerting          *alerting.Config             `json:"alerting,omitempty"`
	Endpoints         []*endpoint.Endpoint         `json:"endpoints,omitempty"`
	ExternalEndpoints []*endpoint.ExternalEndpoint `json:"externalEndpoints,omitempty"`
	Suites            []*suite.Suite               `json:"suites,omitempty"`
}

// ResolveManagedOverlayPath resolves the managed overlay path from either an environment variable or a config path.
func ResolveManagedOverlayPath(configPath string) string {
	if fromEnv := strings.TrimSpace(os.Getenv(GatusManagedOverlayPathEnvVar)); len(fromEnv) > 0 {
		return fromEnv
	}
	if len(configPath) == 0 {
		return ""
	}
	fileInfo, err := os.Stat(configPath)
	if err == nil && fileInfo.IsDir() {
		return filepath.Join(configPath, ".gatus-managed-overlay.json")
	}
	return filepath.Join(filepath.Dir(configPath), ".gatus-managed-overlay.json")
}

// ReadManagedOverlay reads and unmarshals the managed overlay from disk.
// Returns (nil, nil) when no managed overlay exists.
func ReadManagedOverlay(configPath string) (*ManagedOverlay, error) {
	overlayPath := ResolveManagedOverlayPath(configPath)
	if len(overlayPath) == 0 {
		return nil, nil
	}
	overlayBytes, err := os.ReadFile(overlayPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var overlay ManagedOverlay
	if err := json.Unmarshal(overlayBytes, &overlay); err != nil {
		return nil, err
	}
	return &overlay, nil
}

// WriteManagedOverlay writes the managed overlay to disk atomically.
func WriteManagedOverlay(configPath string, overlay *ManagedOverlay) error {
	if overlay == nil {
		return errors.New("managed overlay cannot be nil")
	}
	overlayPath := ResolveManagedOverlayPath(configPath)
	if len(overlayPath) == 0 {
		return errors.New("managed overlay path could not be resolved")
	}
	serialized, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		return err
	}
	serialized = append(serialized, '\n')
	if err := os.MkdirAll(filepath.Dir(overlayPath), 0o755); err != nil {
		return err
	}
	tempFile := fmt.Sprintf("%s.%d.tmp", overlayPath, time.Now().UnixNano())
	if err := os.WriteFile(tempFile, serialized, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tempFile, overlayPath); err != nil {
		_ = os.Remove(tempFile)
		return err
	}
	return nil
}

func applyManagedOverlay(cfg *Config) error {
	overlay, err := ReadManagedOverlay(cfg.configPath)
	if err != nil {
		return err
	}
	if overlay == nil {
		return nil
	}
	overlayHasChanges := false
	if overlay.Alerting != nil {
		cfg.Alerting = overlay.Alerting
		overlayHasChanges = true
	}
	if overlay.Endpoints != nil {
		cfg.Endpoints = overlay.Endpoints
		overlayHasChanges = true
	}
	if overlay.ExternalEndpoints != nil {
		cfg.ExternalEndpoints = overlay.ExternalEndpoints
		overlayHasChanges = true
	}
	if overlay.Suites != nil {
		cfg.Suites = overlay.Suites
		overlayHasChanges = true
	}
	if !overlayHasChanges {
		return nil
	}
	if len(cfg.Endpoints) == 0 && len(cfg.Suites) == 0 {
		return ErrNoEndpointOrSuiteInConfig
	}
	ValidateAlertingConfig(cfg.Alerting, cfg.Endpoints, cfg.ExternalEndpoints)
	if err := ValidateEndpointsConfig(cfg); err != nil {
		return err
	}
	if err := ValidateSuitesConfig(cfg); err != nil {
		return err
	}
	if err := ValidateUniqueKeys(cfg); err != nil {
		return err
	}
	if err := ValidateTunnelingConfig(cfg); err != nil {
		return err
	}
	ValidateAndSetConcurrencyDefaults(cfg)
	if cfg.UI != nil && cfg.Storage != nil {
		cfg.UI.MaximumNumberOfResults = cfg.Storage.MaximumNumberOfResults
	}
	return nil
}

func (config *Config) hasManagedOverlayBeenModified() bool {
	if len(config.managedOverlayPath) == 0 {
		return false
	}
	overlayFileInfo, err := os.Stat(config.managedOverlayPath)
	if err != nil {
		return errors.Is(err, os.ErrNotExist) && !config.lastManagedOverlayModTime.IsZero()
	}
	if config.lastManagedOverlayModTime.IsZero() {
		return true
	}
	return config.lastManagedOverlayModTime.Unix() < overlayFileInfo.ModTime().Unix()
}

func (config *Config) updateManagedOverlayModTime() {
	if len(config.managedOverlayPath) == 0 {
		config.lastManagedOverlayModTime = time.Time{}
		return
	}
	overlayFileInfo, err := os.Stat(config.managedOverlayPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			config.lastManagedOverlayModTime = time.Time{}
		}
		return
	}
	config.lastManagedOverlayModTime = overlayFileInfo.ModTime()
}
