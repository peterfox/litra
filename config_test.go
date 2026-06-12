package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigMissing(t *testing.T) {
	cfg, err := loadConfigFromPath(filepath.Join(t.TempDir(), ".litra"))
	if err != nil {
		t.Fatalf("loadConfigFromPath missing file = %v, want nil error", err)
	}
	if cfg.Profiles == nil {
		t.Error("Profiles map should be initialized, got nil")
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(cfg.Profiles))
	}
}

func TestLoadConfigRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".litra")

	cfg := &Config{
		Profiles: map[string]Profile{
			"work": {Brightness: 70, Temperature: 4500},
			"warm": {Brightness: 30, Temperature: 2800},
		},
	}
	if err := cfg.saveToPath(path); err != nil {
		t.Fatalf("saveToPath: %v", err)
	}

	loaded, err := loadConfigFromPath(path)
	if err != nil {
		t.Fatalf("loadConfigFromPath: %v", err)
	}

	for name, want := range cfg.Profiles {
		got, ok := loaded.Profiles[name]
		if !ok {
			t.Errorf("profile %q missing after roundtrip", name)
			continue
		}
		if got != want {
			t.Errorf("profile %q = %+v, want %+v", name, got, want)
		}
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".litra")
	if err := os.WriteFile(path, []byte("not json"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := loadConfigFromPath(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
