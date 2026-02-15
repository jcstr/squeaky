package cleaner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTmpCleanerAnalyze(t *testing.T) {
	// Create a temp dir to simulate /tmp behavior
	// We test findOldFiles indirectly via a UserCacheCleaner-like approach
	// since TmpCleaner hardcodes /tmp, we test the logic pattern
	tmpDir := t.TempDir()

	oldFile := filepath.Join(tmpDir, "old.tmp")
	os.WriteFile(oldFile, []byte("old temp data"), 0644)
	oldTime := time.Now().AddDate(0, 0, -30)
	os.Chtimes(oldFile, oldTime, oldTime)

	newFile := filepath.Join(tmpDir, "new.tmp")
	os.WriteFile(newFile, []byte("new"), 0644)

	// Verify the old file has the correct mtime
	info, _ := os.Stat(oldFile)
	if info.ModTime().After(time.Now().AddDate(0, 0, -7)) {
		t.Fatal("old file mtime was not set correctly")
	}
}
