package config

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

// WithConfig loads and returns the CLI configuration.
func WithConfig(f func(c *Config)) {
	fmt.Printf("\n\n *** WithConfig *** \n\n")

	WithConfigFrom(DefaultConfigDirectory, f)
}

// WithConfigFrom loads and returns the CLI configuration from a specified location.
func WithConfigFrom(configDir string, f func(c *Config)) {
	fmt.Printf("\n\n *** WithConfigFrom *** \n\n")

	c, err := LoadConfig(configDir)
	if err != nil {
		log.Fatal(err)
	}

	f(c)
}

// WithConfiguration loads and returns the CLI configuration from a specified location.
func WithConfiguration(f func(c *Configuration)) {
	c, err := Configure("")
	if err != nil {
		log.Fatal(err)
	}

	f(c)
}

// TODO: better function name?
func keyGlobalScope(key string) string {
	return fmt.Sprintf("%s.%s", globalScopeIdentifier, key)
}

func getDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.newrelic", home), nil
}
