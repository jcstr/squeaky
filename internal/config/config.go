package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all user-configurable thresholds.
type Config struct {
	PacmanKeep    int      `yaml:"pacman_keep"`
	CacheMaxDays  int      `yaml:"cache_max_days"`
	JournalMaxAge string   `yaml:"journal_max_age"`
	TmpMaxDays    int      `yaml:"tmp_max_days"`
	SkipCleaners  []string `yaml:"skip"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		PacmanKeep:    2,
		CacheMaxDays:  30,
		JournalMaxAge: "2weeks",
		TmpMaxDays:    7,
	}
}

// Load reads a YAML config file and merges it over defaults.
// If the file does not exist, it returns defaults without error.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		configDir, err := os.UserConfigDir()
		if err == nil {
			path = filepath.Join(configDir, "squeaky", "squeaky.yaml")
		}
	}

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return cfg, nil
}
