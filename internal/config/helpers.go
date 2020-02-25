package config

import (
	log "github.com/sirupsen/logrus"
)

// WithConfig loads and returns the CLI configuration.
func WithConfig(f func(c *Config)) {
	WithConfigFrom(DefaultConfigDirectory, f)
}

// WithConfigFrom loads and returns the CLI configuration from a specified location.
func WithConfigFrom(configDir string, f func(c *Config)) {
	c, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot load configuration")
	}

	f(c)
}
