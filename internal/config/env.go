package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	GitlabURL   string `json:"gitlabURL"`
	AccessToken string `json:"accessToken"`
}

func LoadConfig(configPath string) (*Config, error) {
	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	reader := bytes.NewReader(fileBytes)

	config := &Config{}
	err = json.NewDecoder(reader).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	envGitlabURL := os.Getenv("GITLAB_URL")
	envAccessToken := os.Getenv("ACCESS_TOKEN")

	if envGitlabURL != "" {
		config.GitlabURL = envGitlabURL
	}
	if envAccessToken != "" {
		config.AccessToken = envAccessToken
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	var errors []string
	if c.GitlabURL == "" {
		errors = append(errors, "GitlabURL is required")
	}
	if c.AccessToken == "" {
		errors = append(errors, "AccessToken is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}
