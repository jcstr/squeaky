package ui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jcstr/squeaky/internal/cleaner"
)

var (
	banner  = color.New(color.FgCyan)
	title   = color.New(color.FgGreen, color.Bold)
	label   = color.New(color.FgHiWhite, color.Bold)
	would   = color.New(color.FgRed, color.Bold)
	freeing = color.New(color.FgGreen, color.Bold)
	success = color.New(color.FgBlue, color.Bold)
	warning = color.New(color.FgYellow, color.Bold)
	fail    = color.New(color.FgRed, color.Bold)
	dim     = color.New(color.FgHiBlack)
)

// Banner prints the squeaky banner.
func Banner() {
	banner.Println(` 
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
	title.Print("==> ")
	label.Printf("%s\n", name)
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
		fmt.Printf("Would clean ")
		would.Printf("%d", totalItems)
		fmt.Print(" items, freeing ~")
		freeing.Printf("%s\n",cleaner.BytesToHuman(totalBytes))
		warning.Println("Run `squeaky clean` to execute.")
	} else {
		success.Printf("Cleaned %d items, freed %s\n",
			totalItems, cleaner.BytesToHuman(totalBytes))
	}
}
