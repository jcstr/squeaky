package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/jcstr/squeaky/internal/cleaner"
	"github.com/jcstr/squeaky/internal/config"
	"github.com/jcstr/squeaky/internal/ui"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Actually clean your system",
	Long:  "Remove stale packages, caches, logs, and temp files. Prompts for sudo if needed.",
	RunE:  runClean,
}

func init() {
	cleanCmd.Flags().StringSliceVar(&skipList, "skip", nil,
		`cleaners to skip (e.g., --skip "Pacman Cache,Journal Logs")`)
	cleanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false,
		"show individual file details")
	rootCmd.AddCommand(cleanCmd)
}

func runClean(cmd *cobra.Command, args []string) error {
	// If not root, re-exec with sudo (prompts for password)
	if os.Geteuid() != 0 {
		binary, err := exec.LookPath("sudo")
		if err != nil {
			return fmt.Errorf("sudo not found: %w", err)
		}
		sysArgs := append([]string{"sudo"}, os.Args...)
		syscall.Exec(binary, sysArgs, os.Environ())
	}

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
		result, err = c.Clean()

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

	ui.Summary(results, false)
	return nil
}
