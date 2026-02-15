package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// UserCacheCleaner removes old files from ~/.cache.
type UserCacheCleaner struct {
	MaxDays int
}

func NewUserCacheCleaner(maxDays int) *UserCacheCleaner {
	return &UserCacheCleaner{MaxDays: maxDays}
}

func (c *UserCacheCleaner) Name() string {
	return "User Cache"
}

func (c *UserCacheCleaner) cacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, ".cache"), nil
}

func (c *UserCacheCleaner) findOldFiles() ([]string, int64, error) {
	dir, err := c.cacheDir()
	if err != nil {
		return nil, 0, err
	}

	cutoff := time.Now().AddDate(0, 0, -c.MaxDays)
	var files []string
	var totalSize int64

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip permission errors
		}
		if d.IsDir() {
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

func (c *UserCacheCleaner) Analyze() (*Result, error) {
	files, totalSize, err := c.findOldFiles()
	if err != nil {
		return nil, err
	}
	return &Result{
		Name:       c.Name(),
		ItemsFound: len(files),
		SpaceSaved: totalSize,
		Details:    files,
	}, nil
}

func (c *UserCacheCleaner) Clean() (*Result, error) {
	files, _, err := c.findOldFiles()
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
		Name:       c.Name(),
		ItemsFound: removed,
		SpaceSaved: freedBytes,
	}, nil
}
