package sql

import (
	"testing"
	"time"

	"github.com/TwiN/gatus/v5/storage"
	"github.com/TwiN/gatus/v5/storage/store/common"
)

func TestInsertAdminAuditLog(t *testing.T) {
	store, _ := NewStore("sqlite", t.TempDir()+"/TestInsertAdminAuditLog.db", false, storage.DefaultMaximumNumberOfResults, storage.DefaultMaximumNumberOfEvents)
	defer store.Close()

	entry := &common.AdminAuditLogEntry{
		Actor:      "admin@example.com",
		Action:     "create",
		EntityType: "endpoint",
		EntityKey:  "core_api",
		Result:     "success",
		Timestamp:  time.Now().UTC(),
	}
	if err := store.InsertAdminAuditLog(entry); err != nil {
		t.Fatalf("unexpected insert error: %v", err)
	}

	logs, total, err := store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total=1, got %d", total)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(logs))
	}
	if logs[0].Actor != "admin@example.com" {
		t.Fatalf("expected actor admin@example.com, got %s", logs[0].Actor)
	}
	if logs[0].EntityKey != "core_api" {
		t.Fatalf("expected entityKey core_api, got %s", logs[0].EntityKey)
	}
}

func TestInsertAdminAuditLog_Nil(t *testing.T) {
	store, _ := NewStore("sqlite", t.TempDir()+"/TestInsertAdminAuditLog_Nil.db", false, storage.DefaultMaximumNumberOfResults, storage.DefaultMaximumNumberOfEvents)
	defer store.Close()

	if err := store.InsertAdminAuditLog(nil); err != nil {
		t.Fatalf("inserting nil should not return error, got %v", err)
	}
}

func TestGetAdminAuditLogsWithFilters(t *testing.T) {
	store, _ := NewStore("sqlite", t.TempDir()+"/TestGetAdminAuditLogsWithFilters.db", false, storage.DefaultMaximumNumberOfResults, storage.DefaultMaximumNumberOfEvents)
	defer store.Close()

	now := time.Now().UTC()
	entries := []*common.AdminAuditLogEntry{
		{Actor: "alice", Action: "create", EntityType: "endpoint", EntityKey: "a", Result: "success", Timestamp: now},
		{Actor: "bob", Action: "delete", EntityType: "suite", EntityKey: "b", Result: "failure", Error: "not found", Timestamp: now},
		{Actor: "alice", Action: "update", EntityType: "endpoint", EntityKey: "c", Result: "success", Timestamp: now},
	}
	for _, e := range entries {
		if err := store.InsertAdminAuditLog(e); err != nil {
			t.Fatalf("unexpected insert error: %v", err)
		}
	}

	// Filter by actor
	logs, total, err := store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10, Actor: "alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected 2 results for actor=alice, got %d", total)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(logs))
	}

	// Filter by action
	logs, total, err = store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10, Action: "delete"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 result for action=delete, got %d", total)
	}

	// Filter by entityType
	logs, total, err = store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10, EntityType: "suite"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 result for entityType=suite, got %d", total)
	}

	// Filter by result
	logs, total, err = store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10, Result: "failure"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 result for result=failure, got %d", total)
	}
	if logs[0].Error != "not found" {
		t.Fatalf("expected error 'not found', got %s", logs[0].Error)
	}

	// Filter by search
	logs, total, err = store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10, Search: "bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 result for search=bob, got %d", total)
	}
}

func TestGetAdminAuditLogsPagination(t *testing.T) {
	store, _ := NewStore("sqlite", t.TempDir()+"/TestGetAdminAuditLogsPagination.db", false, storage.DefaultMaximumNumberOfResults, storage.DefaultMaximumNumberOfEvents)
	defer store.Close()

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		if err := store.InsertAdminAuditLog(&common.AdminAuditLogEntry{
			Actor:      "admin",
			Action:     "test",
			EntityType: "endpoint",
			Result:     "success",
			Timestamp:  now.Add(time.Duration(i) * time.Second),
		}); err != nil {
			t.Fatalf("unexpected insert error: %v", err)
		}
	}

	// Page 1 with size 2
	logs, total, err := store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total=5, got %d", total)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 results on page 1, got %d", len(logs))
	}

	// Page 3 with size 2 should return 1 result
	logs, total, err = store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 3, PageSize: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total=5, got %d", total)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 result on page 3, got %d", len(logs))
	}

	// Page beyond range should return 0 results
	logs, _, err = store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 10, PageSize: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 0 {
		t.Fatalf("expected 0 results on page 10, got %d", len(logs))
	}
}

func TestDeleteAdminAuditLogsOlderThan(t *testing.T) {
	store, _ := NewStore("sqlite", t.TempDir()+"/TestDeleteAdminAuditLogsOlderThan.db", false, storage.DefaultMaximumNumberOfResults, storage.DefaultMaximumNumberOfEvents)
	defer store.Close()

	now := time.Now().UTC()
	old := now.Add(-48 * time.Hour)
	recent := now.Add(-1 * time.Hour)

	if err := store.InsertAdminAuditLog(&common.AdminAuditLogEntry{
		Actor: "admin", Action: "old-action", EntityType: "endpoint", Result: "success", Timestamp: old,
	}); err != nil {
		t.Fatalf("unexpected insert error: %v", err)
	}
	if err := store.InsertAdminAuditLog(&common.AdminAuditLogEntry{
		Actor: "admin", Action: "recent-action", EntityType: "endpoint", Result: "success", Timestamp: recent,
	}); err != nil {
		t.Fatalf("unexpected insert error: %v", err)
	}

	deleted, err := store.DeleteAdminAuditLogsOlderThan(now.Add(-24 * time.Hour))
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 deleted, got %d", deleted)
	}

	logs, total, err := store.GetAdminAuditLogs(&common.AdminAuditLogQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected 1 remaining log, got %d", total)
	}
	if logs[0].Action != "recent-action" {
		t.Fatalf("expected remaining log to be recent-action, got %s", logs[0].Action)
	}
}
