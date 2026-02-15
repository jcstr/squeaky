package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	skipList []string
	verbose  bool
)

var rootCmd = &cobra.Command{
	Use:   "squeaky",
	Short: "Keep your Arch Linux squeaky clean",
	Long: `Squeaky is an Arch Linux system cleanup tool.
It removes package cache cruft, orphaned packages, old logs,
and stale cache files to keep your system lean.

Run 'squeaky dry' to see what can be cleaned.
Run 'squeaky clean' to actually clean (prompts for sudo).`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default: ~/.config/squeaky/squeaky.yaml)")
}
