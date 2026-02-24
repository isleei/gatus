package api

import (
	"github.com/TwiN/gatus/v5/config"
	"github.com/gofiber/fiber/v2"
)

func TriggerImmediateConfigurationReload(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		config.RequestImmediateReload()
		writeAdminAudit(c, cfg, "reload", "managed-config", "", nil, fiber.Map{"message": "Immediate configuration reload requested."}, nil)
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Immediate configuration reload requested.",
		})
	}
}
