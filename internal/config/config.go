package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Token string `toml:"token"`
}

func Load() (Config, error) {
	var cfg Config

	// Environment variable takes priority
	if token := os.Getenv("BITRISE_TOKEN"); token != "" {
		cfg.Token = token
		return cfg, nil
	}

	// Try config file
	home, err := os.UserHomeDir()
	if err == nil {
		path := filepath.Join(home, ".config", "bitter", "config.toml")
		if _, err := os.Stat(path); err == nil {
			if _, err := toml.DecodeFile(path, &cfg); err != nil {
				return cfg, fmt.Errorf("parsing config file: %w", err)
			}
			if cfg.Token != "" {
				return cfg, nil
			}
		}
	}

	return cfg, fmt.Errorf("BITRISE_TOKEN not set and no token in ~/.config/bitter/config.toml")
}
