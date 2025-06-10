package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadLastMode(t *testing.T) {
	dir := t.TempDir()

	mode := "work"
	if err := SaveLastModeToSettings(dir, mode); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := LoadLastModeFromSettings(dir)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if got != mode {
		t.Fatalf("expected %s got %s", mode, got)
	}

	// verify file exists
	if _, err := os.Stat(filepath.Join(dir, ".todo", "settings.json")); err != nil {
		t.Fatalf("settings file not created: %v", err)
	}
}
