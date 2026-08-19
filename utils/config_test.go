package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configContent := `
settings:
  concurrency: 5
  timeout: 15s

wordpress:
  - name: blog-a
    url: https://blog-a.example.com
    username: metrics-reader
    application_password: test-password
    enabled: true

  - name: blog-b
    url: https://blog-b.example.com
    username: metrics-reader
    application_password: another-password
    enabled: true
`

	configPath := filepath.Join(t.TempDir(), "config.yaml")

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("unable to write test configuration: %v", err)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() returned an error: %v", err)
	}

	if len(config.WordPresses) != 2 {
		t.Fatalf("expected 2 WordPress targets, got %d", len(config.WordPresses))
	}

	if config.WordPresses[0].Name != "blog-a" {
		t.Errorf(
			"expected first target name %q, got %q",
			"blog-a",
			config.WordPresses[0].Name,
		)
	}

	if config.WordPresses[1].Name != "blog-b" {
		t.Errorf(
			"expected second target name %q, got %q",
			"blog-b",
			config.WordPresses[1].Name,
		)
	}

	if config.Settings.Concurrency != 5 {
		t.Errorf(
			"expected concurrency %d, got %d",
			5,
			config.Settings.Concurrency,
		)
	}

	if config.Settings.Timeout != "15s" {
		t.Errorf(
			"expected timeout %q, got %q",
			"15s",
			config.Settings.Timeout,
		)
	}
}
func TestLoadConfigRejectsMissingFile(t *testing.T) {
	_, err := LoadConfig("/AGreatRandomPath/config.yaml")

	if err == nil {
		t.Fatal("expected an error for a missing configuration file")
	}
}

func TestLoadConfigRejectsInvalidYAML(t *testing.T) {
	configContent := `
settings:
  concurrency: [invalidPomme
wordpress:
  - name: blog-abr;;
`

	configPath := filepath.Join(t.TempDir(), "config.yaml")

	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("unable to write test configuration: %v", err)
	}

	if _, err := LoadConfig(configPath); err == nil {
		t.Fatal("expected an error for invalid YAML")
	}
}
