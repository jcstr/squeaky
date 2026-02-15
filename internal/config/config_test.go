package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PacmanKeep != 2 {
		t.Errorf("PacmanKeep = %d, want 2", cfg.PacmanKeep)
	}
	if cfg.CacheMaxDays != 30 {
		t.Errorf("CacheMaxDays = %d, want 30", cfg.CacheMaxDays)
	}
	if cfg.JournalMaxAge != "2weeks" {
		t.Errorf("JournalMaxAge = %q, want %q", cfg.JournalMaxAge, "2weeks")
	}
	if cfg.TmpMaxDays != 7 {
		t.Errorf("TmpMaxDays = %d, want 7", cfg.TmpMaxDays)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PacmanKeep != 2 {
		t.Error("expected defaults when file is missing")
	}
}

func TestLoadValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "squeaky.yaml")

	content := []byte(`
pacman_keep: 3
cache_max_days: 14
journal_max_age: "1week"
tmp_max_days: 3
skip:
  - "Journal Logs"
`)
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PacmanKeep != 3 {
		t.Errorf("PacmanKeep = %d, want 3", cfg.PacmanKeep)
	}
	if cfg.CacheMaxDays != 14 {
		t.Errorf("CacheMaxDays = %d, want 14", cfg.CacheMaxDays)
	}
	if cfg.JournalMaxAge != "1week" {
		t.Errorf("JournalMaxAge = %q, want %q", cfg.JournalMaxAge, "1week")
	}
	if cfg.TmpMaxDays != 3 {
		t.Errorf("TmpMaxDays = %d, want 3", cfg.TmpMaxDays)
	}
	if len(cfg.SkipCleaners) != 1 || cfg.SkipCleaners[0] != "Journal Logs" {
		t.Errorf("SkipCleaners = %v, want [Journal Logs]", cfg.SkipCleaners)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte("pacman_keep:\n  - not\n  a: number"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
