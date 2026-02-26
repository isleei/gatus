package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveManagedOverlayPath_EnvOverride(t *testing.T) {
	t.Setenv(GatusManagedOverlayPathEnvVar, "/custom/overlay.json")
	if got := ResolveManagedOverlayPath("/some/config.yaml"); got != "/custom/overlay.json" {
		t.Fatalf("expected /custom/overlay.json, got %s", got)
	}
}

func TestResolveManagedOverlayPath_EmptyConfigPath(t *testing.T) {
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	if got := ResolveManagedOverlayPath(""); got != "" {
		t.Fatalf("expected empty string, got %s", got)
	}
}

func TestResolveManagedOverlayPath_DirectoryPath(t *testing.T) {
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	dir := t.TempDir()
	expected := filepath.Join(dir, ".gatus-managed-overlay.json")
	if got := ResolveManagedOverlayPath(dir); got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestResolveManagedOverlayPath_FilePath(t *testing.T) {
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configFile, []byte("endpoints: []"), 0o600); err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(dir, ".gatus-managed-overlay.json")
	if got := ResolveManagedOverlayPath(configFile); got != expected {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestReadManagedOverlay_FileNotExist(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	overlay, err := ReadManagedOverlay(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if overlay != nil {
		t.Fatalf("expected nil overlay, got %+v", overlay)
	}
}

func TestReadManagedOverlay_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	overlayPath := filepath.Join(dir, ".gatus-managed-overlay.json")
	data := `{"endpoints":[{"Name":"test","Group":"g","URL":"https://example.org","Interval":30000000000,"Conditions":["[STATUS] == 200"]}]}`
	if err := os.WriteFile(overlayPath, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	overlay, err := ReadManagedOverlay(filepath.Join(dir, "config.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overlay == nil || len(overlay.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %+v", overlay)
	}
	if overlay.Endpoints[0].Name != "test" {
		t.Fatalf("expected endpoint name 'test', got %s", overlay.Endpoints[0].Name)
	}
}

func TestReadManagedOverlay_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	overlayPath := filepath.Join(dir, ".gatus-managed-overlay.json")
	if err := os.WriteFile(overlayPath, []byte("{invalid json}"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := ReadManagedOverlay(filepath.Join(dir, "config.yaml"))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestWriteManagedOverlay_NilOverlayErrors(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	err := WriteManagedOverlay(filepath.Join(dir, "config.yaml"), nil)
	if err == nil {
		t.Fatal("expected error for nil overlay, got nil")
	}
}

func TestWriteManagedOverlay_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	t.Setenv(GatusManagedOverlayPathEnvVar, "")
	overlay := &ManagedOverlay{
		Endpoints: nil,
		Suites:    nil,
	}
	if err := WriteManagedOverlay(configPath, overlay); err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	readBack, err := ReadManagedOverlay(configPath)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if readBack == nil {
		t.Fatal("expected non-nil overlay after write")
	}
}
