package cleaner

import (
	"fmt"
	"os/exec"
)

// Result holds what a cleaner found or did.
type Result struct {
	Name       string
	ItemsFound int
	SpaceSaved int64
	Details    []string
	Skipped    bool
	SkipReason string
	Error      error
}

// Cleaner is the interface every cleanup module implements.
type Cleaner interface {
	// Name returns a human-readable name for this cleaner.
	Name() string

	// Analyze inspects the system and returns what would be cleaned.
	// It must NOT modify anything.
	Analyze() (*Result, error)

	// Clean performs the actual cleanup.
	Clean() (*Result, error)
}

// BytesToHuman converts bytes to a human-readable string.
func BytesToHuman(b int64) string {
	const (
		KiB = 1024
		MiB = KiB * 1024
		GiB = MiB * 1024
	)
	switch {
	case b >= GiB:
		return fmt.Sprintf("%.2f GiB", float64(b)/float64(GiB))
	case b >= MiB:
		return fmt.Sprintf("%.2f MiB", float64(b)/float64(MiB))
	case b >= KiB:
		return fmt.Sprintf("%.2f KiB", float64(b)/float64(KiB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// RunCommand executes a system command and returns its combined output.
// It is a package-level variable so tests can replace it with a fake.
var RunCommand = func(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
