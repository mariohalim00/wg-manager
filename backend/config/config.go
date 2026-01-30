package config

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration.
type Config struct {
	ServerPort string `json:"server_port"`
	// Add other configuration fields here
}

// LoadConfig loads configuration from the specified JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
