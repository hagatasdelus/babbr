package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindConfigFile_CrossPlatform(t *testing.T) {
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalAppData := os.Getenv("APPDATA")

	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
		os.Setenv("APPDATA", originalAppData)
	}()

	tempDir := t.TempDir()

	if runtime.GOOS == "windows" {
		os.Setenv("APPDATA", tempDir)
		os.Unsetenv("XDG_CONFIG_HOME")

		configDir := filepath.Join(tempDir, "babbr")
		err := os.MkdirAll(configDir, 0o700)
		if err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		configFile := filepath.Join(configDir, "config.yaml")
		err = os.WriteFile(configFile, []byte("abbreviations: []"), 0o600)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		found, err := findConfigFile()
		if err != nil {
			t.Fatalf("findConfigFile failed: %v", err)
		}

		expected := filepath.Join(tempDir, "babbr", "config.yaml")
		if found != expected {
			t.Errorf("Expected config file path %s, got %s", expected, found)
		}
	} else {
		os.Setenv("XDG_CONFIG_HOME", tempDir)
		os.Unsetenv("APPDATA")

		configDir := filepath.Join(tempDir, "babbr")
		err := os.MkdirAll(configDir, 0o700)
		if err != nil {
			t.Fatalf("Failed to create config directory: %v", err)
		}

		configFile := filepath.Join(configDir, "config.yaml")
		err = os.WriteFile(configFile, []byte("abbreviations: []"), 0o600)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}

		found, err := findConfigFile()
		if err != nil {
			t.Fatalf("findConfigFile failed: %v", err)
		}

		expected := filepath.Join(tempDir, "babbr", "config.yaml")
		if found != expected {
			t.Errorf("Expected config file path %s, got %s", expected, found)
		}

		os.Unsetenv("XDG_CONFIG_HOME")

		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("Failed to get home directory: %v", err)
		}

		fallbackConfigDir := filepath.Join(homeDir, ".config", "babbr")
		err = os.MkdirAll(fallbackConfigDir, 0o700)
		if err != nil {
			t.Fatalf("Failed to create fallback config directory: %v", err)
		}

		fallbackConfigFile := filepath.Join(fallbackConfigDir, "config.yaml")
		err = os.WriteFile(fallbackConfigFile, []byte("abbreviations: []"), 0o600)
		if err != nil {
			t.Fatalf("Failed to create fallback config file: %v", err)
		}

		found, err = findConfigFile()
		if err != nil {
			t.Fatalf("findConfigFile with fallback failed: %v", err)
		}

		expected = fallbackConfigFile
		if found != expected {
			t.Errorf("Expected fallback config file path %s, got %s", expected, found)
		}

		os.Remove(fallbackConfigFile)
		os.Remove(fallbackConfigDir)
	}
}

func TestFindConfigFile_MissingConfig(t *testing.T) {
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalAppData := os.Getenv("APPDATA")

	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
		os.Setenv("APPDATA", originalAppData)
	}()

	tempDir := t.TempDir()

	if runtime.GOOS == "windows" {
		os.Setenv("APPDATA", tempDir)
		os.Unsetenv("XDG_CONFIG_HOME")
	} else {
		os.Setenv("XDG_CONFIG_HOME", tempDir)
		os.Unsetenv("APPDATA")
	}

	_, err := findConfigFile()
	if err == nil {
		t.Error("Expected findConfigFile to fail when config file doesn't exist")
	}
}
