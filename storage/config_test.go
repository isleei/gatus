package storage

import (
	"testing"
	"time"
)

func TestConfig_ValidateAndSetDefaults_AdminAuditMaxAge(t *testing.T) {
	cfg := &Config{Type: TypeMemory, AdminAuditMaxAge: 720 * time.Hour}
	if err := cfg.ValidateAndSetDefaults(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AdminAuditMaxAge != 720*time.Hour {
		t.Fatalf("expected max age preserved, got %v", cfg.AdminAuditMaxAge)
	}
	cfgNeg := &Config{Type: TypeMemory, AdminAuditMaxAge: -time.Hour}
	if err := cfgNeg.ValidateAndSetDefaults(); err == nil {
		t.Fatal("expected error for negative admin-audit-max-age")
	}
}
