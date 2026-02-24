package common

import "time"

const (
	// DefaultAdminAuditPage is the default page for admin audit log queries.
	DefaultAdminAuditPage = 1

	// DefaultAdminAuditPageSize is the default page size for admin audit log queries.
	DefaultAdminAuditPageSize = 50

	// MaxAdminAuditPageSize is the hard limit for admin audit log page size.
	MaxAdminAuditPageSize = 200
)

// AdminAuditLogEntry represents one admin operation audit log.
type AdminAuditLogEntry struct {
	ID         int64     `json:"id"`
	Actor      string    `json:"actor"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType"`
	EntityKey  string    `json:"entityKey,omitempty"`
	Result     string    `json:"result"` // "success" or "failure"
	Error      string    `json:"error,omitempty"`
	Before     string    `json:"before,omitempty"`
	After      string    `json:"after,omitempty"`
	RequestID  string    `json:"requestID,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// AdminAuditLogQuery defines filters and pagination for querying admin audit logs.
type AdminAuditLogQuery struct {
	Page       int
	PageSize   int
	Actor      string
	Action     string
	EntityType string
	Result     string
	Search     string
	From       *time.Time
	To         *time.Time
}

// Normalize applies sane defaults and limits for pagination.
func (q *AdminAuditLogQuery) Normalize() {
	if q.Page <= 0 {
		q.Page = DefaultAdminAuditPage
	}
	if q.PageSize <= 0 {
		q.PageSize = DefaultAdminAuditPageSize
	}
	if q.PageSize > MaxAdminAuditPageSize {
		q.PageSize = MaxAdminAuditPageSize
	}
}
