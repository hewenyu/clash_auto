package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// FilterRules defines the rules for filtering proxies.
type FilterRules struct {
	IncludeKeywords []string `yaml:"include_keywords"`
}

// Config holds the application configuration.
type Config struct {
	Subscriptions []string    `yaml:"subscriptions"`
	FilterRules   FilterRules `yaml:"filter_rules"`
	TemplatePath  string      `yaml:"template_path"`
	OutputPath    string      `yaml:"output_path"`
}

// LoadConfig reads and parses the configuration file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
