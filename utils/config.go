package utils

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Settings    Settings `yaml:"settings"`
	WordPresses []Target `yaml:"wordpress"`
}

func (c *Config) Validate() error {
	if len(c.WordPresses) == 0 {
		return fmt.Errorf("at least one wordpress target is required")
	}

	if c.Settings.Concurrency <= 0 {
		c.Settings.Concurrency = 10
	}

	if c.Settings.Timeout == "" {
		c.Settings.Timeout = "10s"
	}

	if _, err := time.ParseDuration(c.Settings.Timeout); err != nil {
		return fmt.Errorf("invalid settings.timeout %q: %w", c.Settings.Timeout, err)
	}

	names := make(map[string]struct{})
	urls := make(map[string]struct{})

	for index := range c.WordPresses {
		target := &c.WordPresses[index]

		if err := target.Validate(); err != nil {
			return fmt.Errorf("wordpress[%d]: %w", index, err)
		}

		if _, exists := names[target.Name]; exists {
			return fmt.Errorf("duplicate wordpress target name %q", target.Name)
		}
		names[target.Name] = struct{}{}

		if _, exists := urls[target.URL]; exists {
			return fmt.Errorf("duplicate wordpress target URL %q", target.URL)
		}
		urls[target.URL] = struct{}{}
	}

	return nil
}

type Settings struct {
	Concurrency int    `yaml:"concurrency"`
	Timeout     string `yaml:"timeout"`
}

type Target struct {
	Name                string `yaml:"name"`
	URL                 string `yaml:"url"`
	Username            string `yaml:"username"`
	ApplicationPassword string `yaml:"application_password"`
	Enabled             *bool  `yaml:"enabled"`
	Timeout             string `yaml:"timeout"`
}

func (t *Target) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}

	parsedURL, err := url.Parse(t.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL host is required")
	}

	if t.Username == "" {
		return fmt.Errorf("username is required")
	}

	if t.ApplicationPassword == "" {
		return fmt.Errorf("application_password is required")
	}

	if t.Timeout != "" {
		if _, err := time.ParseDuration(t.Timeout); err != nil {
			return fmt.Errorf("invalid timeout %q: %w", t.Timeout, err)
		}
	}

	return nil
}

func (t *Target) IsEnabled() bool {
	return t.Enabled == nil || *t.Enabled
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read configuration file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse configuration file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	return &config, nil
}
