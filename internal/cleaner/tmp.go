package cleaner

import (
	"os"
	"path/filepath"
	"time"
)

// TmpCleaner removes old files from /tmp.
type TmpCleaner struct {
	MaxDays int
}

func NewTmpCleaner(maxDays int) *TmpCleaner {
	return &TmpCleaner{MaxDays: maxDays}
}

func (t *TmpCleaner) Name() string {
	return "Temp Files"
}

func (t *TmpCleaner) findOldFiles() ([]string, int64, error) {
	cutoff := time.Now().AddDate(0, 0, -t.MaxDays)
	var files []string
	var totalSize int64

	err := filepath.WalkDir("/tmp", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip permission errors
		}
		if path == "/tmp" || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			files = append(files, path)
			totalSize += info.Size()
		}
		return nil
	})

	return files, totalSize, err
}

func (t *TmpCleaner) Analyze() (*Result, error) {
	files, totalSize, err := t.findOldFiles()
	if err != nil {
		return nil, err
	}
	return &Result{
		Name:       t.Name(),
		ItemsFound: len(files),
		SpaceSaved: totalSize,
		Details:    files,
	}, nil
}

func (t *TmpCleaner) Clean() (*Result, error) {
	files, _, err := t.findOldFiles()
	if err != nil {
		return nil, err
	}

	removed := 0
	var freedBytes int64
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		size := info.Size()
		if err := os.Remove(f); err != nil {
			continue
		}
		removed++
		freedBytes += size
	}

	return &Result{
		Name:       t.Name(),
		ItemsFound: removed,
		SpaceSaved: freedBytes,
	}, nil
}
