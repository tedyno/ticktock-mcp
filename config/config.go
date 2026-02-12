package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey      string `json:"api_key"`
	WorkspaceID string `json:"workspace_id,omitempty"`
}

const configDir = "ticktock-mcp"
const configFile = "config.json"

// Load reads configuration from environment variables with fallback to config file.
// Priority: CLOCKIFY_API_KEY env > config file api_key
// Optional: CLOCKIFY_WORKSPACE_ID env > config file workspace_id
func Load() (*Config, error) {
	cfg := &Config{}

	// Try config file first as base
	if fileCfg, err := loadFromFile(); err == nil {
		cfg = fileCfg
	}

	// Env variables override config file
	if key := os.Getenv("CLOCKIFY_API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if wsID := os.Getenv("CLOCKIFY_WORKSPACE_ID"); wsID != "" {
		cfg.WorkspaceID = wsID
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("CLOCKIFY_API_KEY not set (use env variable or ~/.config/%s/%s)", configDir, configFile)
	}

	return cfg, nil
}

func loadFromFile() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".config", configDir, configFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file %s: %w", path, err)
	}

	return &cfg, nil
}
