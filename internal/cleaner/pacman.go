package cleaner

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// PacmanCleaner removes old packages from the pacman cache using paccache.
type PacmanCleaner struct {
	Keep int
}

func NewPacmanCleaner(keep int) *PacmanCleaner {
	return &PacmanCleaner{Keep: keep}
}

func (p *PacmanCleaner) Name() string {
	return "Pacman Cache"
}

func (p *PacmanCleaner) Analyze() (*Result, error) {
	out, err := RunCommand("paccache", "-d", "-k", strconv.Itoa(p.Keep))
	if err != nil {
		return nil, fmt.Errorf("paccache dry-run: %w: %s", err, out)
	}
	return parsePaccacheOutput(p.Name(), out), nil
}

func (p *PacmanCleaner) Clean() (*Result, error) {
	out, err := RunCommand("sudo", "paccache", "-r", "-k", strconv.Itoa(p.Keep))
	if err != nil {
		if strings.Contains(out, "ermission denied") {
			return &Result{
				Name:       p.Name(),
				Skipped:    true,
				SkipReason: "requires root (run with sudo)",
			}, nil
		}
		return nil, fmt.Errorf("paccache remove: %w: %s", err, out)
	}
	return parsePaccacheOutput(p.Name(), out), nil
}

var paccacheRe = regexp.MustCompile(
	`(\d+)\s+candidates.*?([\d.]+)\s+(GiB|MiB|KiB|B)`,
)

func parsePaccacheOutput(name, output string) *Result {
	r := &Result{Name: name}
	matches := paccacheRe.FindStringSubmatch(output)
	if matches == nil {
		return r
	}
	r.ItemsFound, _ = strconv.Atoi(matches[1])
	size, _ := strconv.ParseFloat(matches[2], 64)
	switch matches[3] {
	case "GiB":
		r.SpaceSaved = int64(size * 1024 * 1024 * 1024)
	case "MiB":
		r.SpaceSaved = int64(size * 1024 * 1024)
	case "KiB":
		r.SpaceSaved = int64(size * 1024)
	default:
		r.SpaceSaved = int64(size)
	}
	return r
}
