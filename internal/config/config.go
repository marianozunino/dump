package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Constants for default paths
const (
	DefaultUploadPath = "./uploads"
	DefaultConfigPath = "./config/config.json"
)

// Constants for upload settings
const (
	MaxUploadSize   = 100 * 1024 * 1024 // 100MB
	DefaultIDLength = 4
)

// Config represents the application configuration
type Config struct {
	MinAge        int     `json:"min_age_days"`       // Minimum retention in days
	MaxAge        int     `json:"max_age_days"`       // Maximum retention in days
	MaxSize       float64 `json:"max_size_mib"`       // Maximum file size in MiB
	UploadPath    string  `json:"upload_path"`        // Path to uploaded files
	CheckInterval int     `json:"check_interval_min"` // How often to check for expired files (minutes)
	Enabled       bool    `json:"enabled"`            // Whether expiration is enabled
	BaseURL       string  `json:"base_url"`           // Base URL for links
}

// DefaultConfig provides default config values
var DefaultConfig = Config{
	MinAge:        30,    // 30 days
	MaxAge:        365,   // 1 year
	MaxSize:       512.0, // 512 MiB
	UploadPath:    DefaultUploadPath,
	CheckInterval: 60, // Check once per hour
	Enabled:       true,
	BaseURL:       "http://localhost:8080/", // Change to your domain in production
}

// SetupDefaultConfig creates a default configuration file if none exists
func SetupDefaultConfig() error {
	if _, err := os.Stat(DefaultConfigPath); err == nil {
		return nil
	}

	data, err := json.MarshalIndent(DefaultConfig, "", "  ")
	if err != nil {
		return err
	}

	configDir := filepath.Dir(DefaultConfigPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(DefaultConfigPath, data, 0o644)
}

// LoadConfig loads a configuration from file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}
