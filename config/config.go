package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	MaxAnswers int    `json:"max_answers"`
	Files      []File `json:"files"`
}

// File represents a file source in the configuration
type File struct {
	Type string `json:"type"`
	Path string `json:"path,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Function to load configuration from config.yaml
func LoadConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
