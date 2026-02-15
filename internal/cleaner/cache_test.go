package cleaner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUserCacheAnalyze(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cacheDir := filepath.Join(tmpDir, ".cache")
	os.MkdirAll(cacheDir, 0755)

	// Create an old file (60 days ago)
	oldFile := filepath.Join(cacheDir, "old.txt")
	os.WriteFile(oldFile, []byte("old content here"), 0644)
	oldTime := time.Now().AddDate(0, 0, -60)
	os.Chtimes(oldFile, oldTime, oldTime)

	// Create a new file
	newFile := filepath.Join(cacheDir, "new.txt")
	os.WriteFile(newFile, []byte("new"), 0644)

	c := NewUserCacheCleaner(30)
	result, err := c.Analyze()
	if err != nil {
		t.Fatal(err)
	}
	if result.ItemsFound != 1 {
		t.Errorf("ItemsFound = %d, want 1", result.ItemsFound)
	}
	if result.SpaceSaved != 16 {
		t.Errorf("SpaceSaved = %d, want 16", result.SpaceSaved)
	}
}

func TestUserCacheClean(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cacheDir := filepath.Join(tmpDir, ".cache")
	os.MkdirAll(cacheDir, 0755)

	oldFile := filepath.Join(cacheDir, "stale.dat")
	os.WriteFile(oldFile, []byte("delete me"), 0644)
	oldTime := time.Now().AddDate(0, 0, -60)
	os.Chtimes(oldFile, oldTime, oldTime)

	c := NewUserCacheCleaner(30)
	result, err := c.Clean()
	if err != nil {
		t.Fatal(err)
	}
	if result.ItemsFound != 1 {
		t.Errorf("ItemsFound = %d, want 1", result.ItemsFound)
	}

	// Verify file was deleted
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("old file should have been deleted")
	}
}
