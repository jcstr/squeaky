package cleaner

import (
	"fmt"
	"testing"
)

func TestOrphanAnalyzeNoOrphans(t *testing.T) {
	orig := RunCommand
	defer func() { RunCommand = orig }()

	RunCommand = func(name string, args ...string) (string, error) {
		return "", fmt.Errorf("exit status 1")
	}

	o := NewOrphanCleaner()
	result, err := o.Analyze()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ItemsFound != 0 {
		t.Errorf("ItemsFound = %d, want 0", result.ItemsFound)
	}
}

func TestOrphanAnalyzeWithOrphans(t *testing.T) {
	orig := RunCommand
	defer func() { RunCommand = orig }()

	RunCommand = func(name string, args ...string) (string, error) {
		return "libfoo\nlibbar\nlibold\n", nil
	}

	o := NewOrphanCleaner()
	result, err := o.Analyze()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ItemsFound != 3 {
		t.Errorf("ItemsFound = %d, want 3", result.ItemsFound)
	}
}
