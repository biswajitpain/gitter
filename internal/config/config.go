package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration for the LLM provider.
type Config struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
}

// GetConfigPath returns the path to the configuration file.
// It ensures the parent directory exists.
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	configDir := filepath.Join(home, ".config", "gitter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("could not create config directory: %w", err)
	}
	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads the configuration from the file.
func LoadConfig() (Config, error) {
	var config Config
	path, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty config if file doesn't exist
			return Config{}, nil
		}
		return config, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, fmt.Errorf("could not decode config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to the file.
// It sets file permissions to 0600 for security.
func SaveConfig(config Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("could not create or open config file for writing: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("could not encode config to file: %w", err)
	}

	return nil
}
