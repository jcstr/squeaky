package ui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jcstr/squeaky/internal/cleaner"
)

var (
	title   = color.New(color.FgCyan, color.Bold)
	success = color.New(color.FgGreen)
	warning = color.New(color.FgYellow)
	fail    = color.New(color.FgRed)
	dim     = color.New(color.FgHiBlack)
)

// Banner prints the squeaky banner.
func Banner() {
	title.Println(`                                          
 ▄▄▄▄  ▄▄▄  ▄▄ ▄▄ ▄▄▄▄▄  ▄▄▄  ▄▄ ▄▄ ▄▄ ▄▄ 
███▄▄ ██▀██ ██ ██ ██▄▄  ██▀██ ██▄█▀ ▀███▀ 
▄▄██▀ ▀███▀ ▀███▀ ██▄▄▄ ██▀██ ██ ██   █   
         ▀▀                               `)
	fmt.Println()
}

// DryRunBanner prints the dry-run notice.
func DryRunBanner() {
	warning.Println("DRY RUN: showing what would be cleaned.")
	fmt.Println()
}

// Header prints a section header for a cleaner.
func Header(name string) {
	title.Printf("==> %s\n", name)
}

// PrintResult prints the result of a single cleaner.
func PrintResult(r *cleaner.Result) {
	if r.Skipped {
		dim.Printf("    [SKIP] %s\n", r.SkipReason)
		return
	}
	if r.Error != nil {
		fail.Printf("    [FAIL] %v\n", r.Error)
		return
	}
	if r.ItemsFound == 0 {
		success.Println("    [OK]   nothing to clean")
		return
	}
	success.Printf("    [OK]   %d items, %s\n",
		r.ItemsFound, cleaner.BytesToHuman(r.SpaceSaved))
}

// Summary prints the final tally.
func Summary(results []*cleaner.Result, dryRun bool) {
	var totalItems int
	var totalBytes int64
	for _, r := range results {
		if !r.Skipped && r.Error == nil {
			totalItems += r.ItemsFound
			totalBytes += r.SpaceSaved
		}
	}

	fmt.Println()
	title.Println("--- Summary ---")
	if dryRun {
		fmt.Printf("Would clean %d items, freeing ~%s\n",
			totalItems, cleaner.BytesToHuman(totalBytes))
		warning.Println("Run squeaky clean to execute.")
	} else {
		success.Printf("Cleaned %d items, freed %s\n",
			totalItems, cleaner.BytesToHuman(totalBytes))
	}
}
