package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Abbreviations []Abbreviation `yaml:"abbreviations"`
}

type Abbreviation struct {
	Name    string               `yaml:"name,omitempty"`
	Abbr    string               `yaml:"abbr,omitempty"`
	Snippet string               `yaml:"snippet"`
	Options *AbbreviationOptions `yaml:"options,omitempty"`
}

type AbbreviationOptions struct {
	Position  string `yaml:"position,omitempty"`
	Command   string `yaml:"command,omitempty"`
	Regex     string `yaml:"regex,omitempty"`
	SetCursor bool   `yaml:"set_cursor,omitempty"`
	Evaluate  bool   `yaml:"evaluate,omitempty"`
	Condition string `yaml:"condition,omitempty"`
}

func LoadConfig() (*Config, error) {
	configPath, err := findConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to find config file: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func findConfigFile() (string, error) {
	var configHome string

	if runtime.GOOS == "windows" {
		configHome = os.Getenv("APPDATA")
	} else {
		configHome = os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			configHome = filepath.Join(homeDir, ".config")
		}
	}

	if configHome == "" {
		return "", fmt.Errorf("could not determine config directory")
	}

	dir := filepath.Join(configHome, "babbr")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	userConfigPath := filepath.Join(dir, "config.yaml")
	if _, err := os.Stat(userConfigPath); err == nil {
		return userConfigPath, nil
	}

	return "", fmt.Errorf("config file not found in config directory %s", dir)
}
