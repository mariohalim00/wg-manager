package config

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration.
type Config struct {
	ServerPort         string `json:"server_port"`
	InterfaceName      string `json:"interface_name"`
	StoragePath        string `json:"storage_path"`
	ServerEndpoint     string `json:"server_endpoint"` // e.g. "vpn.example.com:51820"
	ServerPubKey       string `json:"server_pubkey"`
	CORSAllowedOrigins string `json:"cors_allowed_origins"`
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

	// Environment variable overrides
	if envPort := os.Getenv("WG_SERVER_PORT"); envPort != "" {
		cfg.ServerPort = envPort
	}
	if envInterface := os.Getenv("WG_INTERFACE_NAME"); envInterface != "" {
		cfg.InterfaceName = envInterface
	}
	if envStorage := os.Getenv("WG_STORAGE_PATH"); envStorage != "" {
		cfg.StoragePath = envStorage
	}
	if envEndpoint := os.Getenv("WG_SERVER_ENDPOINT"); envEndpoint != "" {
		cfg.ServerEndpoint = envEndpoint
	}
	if envPubKey := os.Getenv("WG_SERVER_PUBKEY"); envPubKey != "" {
		cfg.ServerPubKey = envPubKey
	}
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		cfg.CORSAllowedOrigins = envOrigins
	}

	return &cfg, nil
}
