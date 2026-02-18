package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/TwiN/gatus/v5/alerting"
	"github.com/TwiN/gatus/v5/config"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/gatus/v5/config/suite"
	"github.com/gofiber/fiber/v2"
)

type ManagedConfigPayload struct {
	Alerting          *alerting.Config             `json:"alerting,omitempty"`
	Endpoints         []*endpoint.Endpoint         `json:"endpoints"`
	ExternalEndpoints []*endpoint.ExternalEndpoint `json:"externalEndpoints"`
	Suites            []*suite.Suite               `json:"suites"`
}

type ManagedConfigResponse struct {
	OverlayPath       string                       `json:"overlayPath"`
	Alerting          *alerting.Config             `json:"alerting,omitempty"`
	Endpoints         []*endpoint.Endpoint         `json:"endpoints"`
	ExternalEndpoints []*endpoint.ExternalEndpoint `json:"externalEndpoints"`
	Suites            []*suite.Suite               `json:"suites"`
	Note              string                       `json:"note,omitempty"`
}

func loadManagedCandidate(cfg *config.Config) (*config.Config, error) {
	if cfg == nil {
		return nil, errors.New("nil configuration")
	}
	if len(cfg.LoadedConfigPath()) == 0 {
		candidate := *cfg
		return &candidate, nil
	}
	return config.LoadConfiguration(cfg.LoadedConfigPath())
}

func persistManagedCandidate(cfg *config.Config, candidate *config.Config) error {
	return persistManagedCandidateWithAlerting(cfg, candidate, false)
}

func persistManagedCandidateWithAlerting(cfg *config.Config, candidate *config.Config, persistAlerting bool) error {
	if cfg == nil || candidate == nil {
		return errors.New("configuration cannot be nil")
	}
	if len(cfg.LoadedConfigPath()) == 0 {
		return fmt.Errorf("managed configuration cannot be persisted without a loaded configuration path")
	}
	existingOverlay, err := config.ReadManagedOverlay(cfg.LoadedConfigPath())
	if err != nil {
		return err
	}
	overlay := &config.ManagedOverlay{
		Endpoints:         candidate.Endpoints,
		ExternalEndpoints: candidate.ExternalEndpoints,
		Suites:            candidate.Suites,
	}
	if existingOverlay != nil {
		overlay.Alerting = existingOverlay.Alerting
	}
	if persistAlerting {
		overlay.Alerting = candidate.Alerting
	}
	return config.WriteManagedOverlay(cfg.LoadedConfigPath(), overlay)
}

func validateManagedPayload(candidate *config.Config) error {
	if len(candidate.Endpoints) == 0 && len(candidate.Suites) == 0 {
		return config.ErrNoEndpointOrSuiteInConfig
	}
	config.ValidateAlertingConfig(candidate.Alerting, candidate.Endpoints, candidate.ExternalEndpoints)
	if err := config.ValidateSecurityConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateEndpointsConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateWebConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateUIConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateMaintenanceConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateStorageConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateRemoteConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateConnectivityConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateSuitesConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateTunnelingConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateAnnouncementsConfig(candidate); err != nil {
		return err
	}
	if err := config.ValidateUniqueKeys(candidate); err != nil {
		return err
	}
	config.ValidateAndSetConcurrencyDefaults(candidate)
	if candidate.UI != nil && candidate.Storage != nil {
		candidate.UI.MaximumNumberOfResults = candidate.Storage.MaximumNumberOfResults
	}
	return nil
}

func GetManagedConfiguration(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentConfig, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(ManagedConfigResponse{
			OverlayPath:       currentConfig.ManagedOverlayPath(),
			Alerting:          currentConfig.Alerting,
			Endpoints:         currentConfig.Endpoints,
			ExternalEndpoints: currentConfig.ExternalEndpoints,
			Suites:            currentConfig.Suites,
			Note:              "Saving managed configuration writes an overlay file and is automatically reloaded within a few seconds.",
		})
	}
}

func PutManagedConfiguration(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload ManagedConfigPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid payload: " + err.Error(),
			})
		}
		candidate, err := loadManagedCandidate(cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		candidate.Alerting = payload.Alerting
		candidate.Endpoints = payload.Endpoints
		candidate.ExternalEndpoints = payload.ExternalEndpoints
		candidate.Suites = payload.Suites
		if err := validateManagedPayload(candidate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err := persistManagedCandidateWithAlerting(cfg, candidate, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"overlayPath": config.ResolveManagedOverlayPath(cfg.LoadedConfigPath()),
			"message":     "Managed configuration saved. Gatus will apply it automatically within a few seconds.",
		})
	}
}

func DeleteManagedConfiguration(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		overlayPath := config.ResolveManagedOverlayPath(cfg.LoadedConfigPath())
		if err := os.Remove(overlayPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
