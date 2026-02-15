package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jcstr/squeaky/internal/cleaner"
	"github.com/jcstr/squeaky/internal/config"
	"github.com/jcstr/squeaky/internal/ui"
)

var dryCmd = &cobra.Command{
	Use:   "dry",
	Short: "Show what would be cleaned (dry-run)",
	Long:  "Analyze your system and show what would be cleaned without making any changes.",
	RunE:  runDry,
}

func init() {
	dryCmd.Flags().StringSliceVar(&skipList, "skip", nil,
		`cleaners to skip (e.g., --skip "Pacman Cache,Journal Logs")`)
	dryCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"show individual file details")
	rootCmd.AddCommand(dryCmd)
}

func runDry(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	skipSet := make(map[string]bool)
	for _, s := range cfg.SkipCleaners {
		skipSet[s] = true
	}
	for _, s := range skipList {
		skipSet[s] = true
	}

	allCleaners := []cleaner.Cleaner{
		cleaner.NewPacmanCleaner(cfg.PacmanKeep),
		cleaner.NewOrphanCleaner(),
		cleaner.NewUserCacheCleaner(cfg.CacheMaxDays),
		cleaner.NewJournalCleaner(cfg.JournalMaxAge),
		cleaner.NewTmpCleaner(cfg.TmpMaxDays),
	}

	ui.Banner()
	ui.DryRunBanner()

	var results []*cleaner.Result

	for _, c := range allCleaners {
		if skipSet[c.Name()] {
			results = append(results, &cleaner.Result{
				Name:       c.Name(),
				Skipped:    true,
				SkipReason: "skipped by user",
			})
			continue
		}

		ui.Header(c.Name())

		var result *cleaner.Result
		result, err = c.Analyze()

		if err != nil {
			fmt.Fprintf(os.Stderr, "    error: %v\n", err)
			results = append(results, &cleaner.Result{
				Name:  c.Name(),
				Error: err,
			})
			continue
		}

		ui.PrintResult(result)
		if verbose && len(result.Details) > 0 {
			for _, d := range result.Details {
				fmt.Printf("      %s\n", d)
			}
		}

		results = append(results, result)
	}

	ui.Summary(results, true)
	return nil
}
