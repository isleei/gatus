package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/TwiN/gatus/v5/config"
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

// DeleteAdminAuditLogsOlderThan deletes admin audit logs older than the given max age (query: maxAge, e.g. 720h).
// Also accepts "days" as an integer alternative (e.g. days=30).
func DeleteAdminAuditLogsOlderThan(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var maxAge time.Duration
		var err error
		if raw := strings.TrimSpace(c.Query("maxAge")); len(raw) > 0 {
			maxAge, err = time.ParseDuration(raw)
			if err != nil {
				writeAdminAudit(c, cfg, "cleanup", "audit-log", "", fiber.Map{"maxAge": raw}, nil, err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid maxAge duration: " + err.Error()})
			}
		} else if days := c.QueryInt("days", 0); days > 0 {
			maxAge = time.Duration(days) * 24 * time.Hour
		} else if cfg != nil && cfg.Storage != nil && cfg.Storage.AdminAuditMaxAge > 0 {
			maxAge = cfg.Storage.AdminAuditMaxAge
		} else {
			err = fmt.Errorf("maxAge or days query parameter required (or configure storage.admin-audit-max-age)")
			writeAdminAudit(c, cfg, "cleanup", "audit-log", "", nil, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if maxAge <= 0 {
			err = fmt.Errorf("maxAge must be positive")
			writeAdminAudit(c, cfg, "cleanup", "audit-log", "", fiber.Map{"maxAge": maxAge.String()}, nil, err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		before := time.Now().Add(-maxAge)
		deleted, err := store.Get().DeleteAdminAuditLogsOlderThan(before)
		if err != nil {
			writeAdminAudit(c, cfg, "cleanup", "audit-log", "", fiber.Map{"before": before}, nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		result := fiber.Map{"deleted": deleted, "before": before.Format(time.RFC3339), "maxAge": maxAge.String()}
		writeAdminAudit(c, cfg, "cleanup", "audit-log", "", fiber.Map{"maxAge": maxAge.String()}, result, nil)
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
