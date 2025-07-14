package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application's configuration.
type Config struct {
	Scoring         Scoring `yaml:"scoring"`
	Events          []Event `yaml:"events"`
	DrupalPublisher struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"drupal_publisher"`
}

// Scoring holds the scoring configuration.
type Scoring struct {
	MaxEventsToCount int `yaml:"max_events_to_count"`
	MaxPoints        int `yaml:"max_points"`
}

// Event represents a single event in the league.
type Event struct {
	Name        string      `yaml:"name"`
	ClosingDate ClosingDate `yaml:"closing_date"`
}

// Load reads a configuration file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
