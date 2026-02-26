package common

import "testing"

func TestAdminAuditLogQuery_Normalize_Defaults(t *testing.T) {
	q := &AdminAuditLogQuery{}
	q.Normalize()
	if q.Page != DefaultAdminAuditPage {
		t.Fatalf("expected Page=%d, got %d", DefaultAdminAuditPage, q.Page)
	}
	if q.PageSize != DefaultAdminAuditPageSize {
		t.Fatalf("expected PageSize=%d, got %d", DefaultAdminAuditPageSize, q.PageSize)
	}
}

func TestAdminAuditLogQuery_Normalize_NegativeValues(t *testing.T) {
	q := &AdminAuditLogQuery{Page: -5, PageSize: -10}
	q.Normalize()
	if q.Page != DefaultAdminAuditPage {
		t.Fatalf("expected Page=%d for negative input, got %d", DefaultAdminAuditPage, q.Page)
	}
	if q.PageSize != DefaultAdminAuditPageSize {
		t.Fatalf("expected PageSize=%d for negative input, got %d", DefaultAdminAuditPageSize, q.PageSize)
	}
}

func TestAdminAuditLogQuery_Normalize_LargePageSize(t *testing.T) {
	q := &AdminAuditLogQuery{Page: 1, PageSize: 999}
	q.Normalize()
	if q.PageSize != MaxAdminAuditPageSize {
		t.Fatalf("expected PageSize capped to %d, got %d", MaxAdminAuditPageSize, q.PageSize)
	}
}

func TestAdminAuditLogQuery_Normalize_ValidValues(t *testing.T) {
	q := &AdminAuditLogQuery{Page: 3, PageSize: 25}
	q.Normalize()
	if q.Page != 3 {
		t.Fatalf("expected Page=3, got %d", q.Page)
	}
	if q.PageSize != 25 {
		t.Fatalf("expected PageSize=25, got %d", q.PageSize)
	}
}
