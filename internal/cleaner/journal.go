package cleaner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// JournalCleaner vacuums old systemd journal logs.
type JournalCleaner struct {
	MaxAge string
}

func NewJournalCleaner(maxAge string) *JournalCleaner {
	return &JournalCleaner{MaxAge: maxAge}
}

func (j *JournalCleaner) Name() string {
	return "Journal Logs"
}

var journalSizeRe = regexp.MustCompile(`([\d.]+)([KMGT]?)`)

func parseJournalSize(output string) int64 {
	matches := journalSizeRe.FindStringSubmatch(output)
	if matches == nil {
		return 0
	}
	size, _ := strconv.ParseFloat(matches[1], 64)
	switch matches[2] {
	case "G":
		return int64(size * 1024 * 1024 * 1024)
	case "M":
		return int64(size * 1024 * 1024)
	case "K":
		return int64(size * 1024)
	default:
		return int64(size)
	}
}

func (j *JournalCleaner) Analyze() (*Result, error) {
	out, err := RunCommand("journalctl", "--disk-usage")
	if err != nil {
		return nil, fmt.Errorf("journalctl --disk-usage: %w: %s", err, out)
	}
	size := parseJournalSize(out)
	return &Result{
		Name:       j.Name(),
		ItemsFound: 1,
		SpaceSaved: size,
		Details:    []string{strings.TrimSpace(out)},
	}, nil
}

func (j *JournalCleaner) Clean() (*Result, error) {
	beforeOut, _ := RunCommand("journalctl", "--disk-usage")
	beforeSize := parseJournalSize(beforeOut)

	out, err := RunCommand("sudo", "journalctl", "--vacuum-time="+j.MaxAge)
	if err != nil {
		if strings.Contains(out, "ermission denied") {
			return &Result{
				Name:       j.Name(),
				Skipped:    true,
				SkipReason: "requires root (run with sudo)",
			}, nil
		}
		return nil, fmt.Errorf("journalctl vacuum: %w: %s", err, out)
	}

	afterOut, _ := RunCommand("journalctl", "--disk-usage")
	afterSize := parseJournalSize(afterOut)
	freed := beforeSize - afterSize
	if freed < 0 {
		freed = 0
	}

	return &Result{
		Name:       j.Name(),
		ItemsFound: 1,
		SpaceSaved: freed,
		Details:    []string{strings.TrimSpace(out)},
	}, nil
}
