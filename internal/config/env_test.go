package config

import "testing"

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		configPath  string
		injectedEnv map[string]string
		expected    Config
	}{
		{
			name:        "Valid config",
			configPath:  "testdata/valid-config.json",
			injectedEnv: map[string]string{},
			expected: Config{
				GitlabURL:   "https://gitlab.example.com",
				AccessToken: "secret-token",
			},
		},
		{
			name:       "Environment variable override",
			configPath: "testdata/valid-config.json",
			injectedEnv: map[string]string{
				"GITLAB_URL":   "https://super.gitlab.example.com",
				"ACCESS_TOKEN": "SECRET",
			},
			expected: Config{
				GitlabURL:   "https://super.gitlab.example.com",
				AccessToken: "SECRET",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.injectedEnv {
				t.Setenv(key, value)
			}

			config, err := LoadConfig(tc.configPath)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if config == nil {
				t.Fatal("Expected config to be non-nil")
			}

			if config.GitlabURL != tc.expected.GitlabURL {
				t.Errorf("Expected GitlabURL '%s', got '%s'", tc.expected.GitlabURL, config.GitlabURL)
			}

			if config.AccessToken != tc.expected.AccessToken {
				t.Errorf("Expected AccessToken '%s', got '%s'", tc.expected.AccessToken, config.AccessToken)
			}
		})
	}

	t.Run("Invalid config file", func(t *testing.T) {
		_, err := LoadConfig("testdata/invalid-config.json")
		if err == nil {
			t.Fatal("Expected error for invalid config file, got nil")
		}
	})

	t.Run("Non existent config file", func(t *testing.T) {
		_, err := LoadConfig("testdata/non-existent.json")
		if err == nil {
			t.Fatal("Expected error for non-existent config file, got nil")
		}
	})

	t.Run("Empty config file", func(t *testing.T) {
		_, err := LoadConfig("testdata/empty-config.json")
		if err == nil {
			t.Fatal("Expected error for empty config file, got nil")
		}
	})
}
