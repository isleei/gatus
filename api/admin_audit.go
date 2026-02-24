package api

import (
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/storage/store"
	"github.com/TwiN/gatus/v5/storage/store/common"
	"github.com/gofiber/fiber/v2"
)

type AdminAuditLogListResponse struct {
	Items    []*common.AdminAuditLogEntry `json:"items"`
	Total    int                          `json:"total"`
	Page     int                          `json:"page"`
	PageSize int                          `json:"pageSize"`
}

func GetAdminAuditLogs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := &common.AdminAuditLogQuery{
			Actor:      c.Query("actor"),
			Action:     c.Query("action"),
			EntityType: c.Query("entityType"),
			Result:     normalizeAdminAuditResultFilter(c.Query("result")),
			Search:     c.Query("q"),
			Page:       c.QueryInt("page", common.DefaultAdminAuditPage),
			PageSize:   c.QueryInt("pageSize", common.DefaultAdminAuditPageSize),
		}
		if from := c.Query("from"); len(from) > 0 {
			if parsed, err := time.Parse(time.RFC3339, from); err == nil {
				query.From = &parsed
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid from timestamp, expected RFC3339"})
			}
		}
		if to := c.Query("to"); len(to) > 0 {
			if parsed, err := time.Parse(time.RFC3339, to); err == nil {
				query.To = &parsed
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid to timestamp, expected RFC3339"})
			}
		}
		query.Normalize()
		items, total, err := store.Get().GetAdminAuditLogs(query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(AdminAuditLogListResponse{
			Items:    items,
			Total:    total,
			Page:     query.Page,
			PageSize: query.PageSize,
		})
	}
}

func normalizeAdminAuditResultFilter(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "success", "failure":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}
