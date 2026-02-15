package cleaner

import (
	"fmt"
	"strings"
)

// OrphanCleaner removes orphaned packages via pacman.
type OrphanCleaner struct{}

func NewOrphanCleaner() *OrphanCleaner {
	return &OrphanCleaner{}
}

func (o *OrphanCleaner) Name() string {
	return "Orphaned Packages"
}

func (o *OrphanCleaner) Analyze() (*Result, error) {
	out, err := RunCommand("pacman", "-Qdtq")
	if err != nil {
		// Exit code 1 with empty output means no orphans
		if strings.TrimSpace(out) == "" {
			return &Result{Name: o.Name(), ItemsFound: 0}, nil
		}
		return nil, fmt.Errorf("pacman -Qdtq: %w: %s", err, out)
	}

	pkgs := strings.Split(strings.TrimSpace(out), "\n")
	if len(pkgs) == 1 && pkgs[0] == "" {
		pkgs = nil
	}

	return &Result{
		Name:       o.Name(),
		ItemsFound: len(pkgs),
		Details:    pkgs,
	}, nil
}

func (o *OrphanCleaner) Clean() (*Result, error) {
	analyzeResult, err := o.Analyze()
	if err != nil {
		return nil, err
	}
	if analyzeResult.ItemsFound == 0 {
		return analyzeResult, nil
	}

	args := append([]string{"pacman", "-Rns", "--noconfirm"}, analyzeResult.Details...)
	out, err := RunCommand("sudo", args...)
	if err != nil {
		if strings.Contains(out, "ermission denied") {
			return &Result{
				Name:       o.Name(),
				Skipped:    true,
				SkipReason: "requires root (run with sudo)",
			}, nil
		}
		return nil, fmt.Errorf("pacman -Rns: %w: %s", err, out)
	}

	return &Result{
		Name:       o.Name(),
		ItemsFound: analyzeResult.ItemsFound,
		Details:    analyzeResult.Details,
	}, nil
}
