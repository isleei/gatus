package api

import (
	"github.com/TwiN/gatus/v5/config"
	"github.com/gofiber/fiber/v2"
)

func TriggerImmediateConfigurationReload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		config.RequestImmediateReload()
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Immediate configuration reload requested.",
		})
	}
}
