package config_test

import (
	"gitter/internal/config"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	// Create a temporary directory for the test to avoid interfering with user's config
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir) // Temporarily set HOME to our temp dir

	// Define a test config
	testConfig := config.Config{
		Provider: "openai",
		APIKey:   "sk-test-key",
	}

	// Test SaveConfig
	err := config.SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Check if the file was created with the correct path
	configPath, err := config.GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}
	expectedPath := filepath.Join(tempDir, ".config", "gitter", "config.json")
	if configPath != expectedPath {
		t.Errorf("GetConfigPath() returned %s, want %s", configPath, expectedPath)
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("os.Stat() on config file failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("config file permissions are %s, want 0600", info.Mode().Perm())
	}

	// Test LoadConfig
	loadedConfig, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Verify the loaded config matches the saved config
	if loadedConfig.Provider != testConfig.Provider {
		t.Errorf("loaded provider is %s, want %s", loadedConfig.Provider, testConfig.Provider)
	}
	if loadedConfig.APIKey != testConfig.APIKey {
		t.Errorf("loaded API key is %s, want %s", loadedConfig.APIKey, testConfig.APIKey)
	}
}

func TestLoadConfig_NotExist(t *testing.T) {
	// Set HOME to a temporary directory where the config does not exist
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// LoadConfig should return an empty config and no error
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() with no existing file failed: %v", err)
	}
	if cfg.Provider != "" || cfg.APIKey != "" {
		t.Errorf("LoadConfig() with no existing file should return an empty config, but got %+v", cfg)
	}
}
