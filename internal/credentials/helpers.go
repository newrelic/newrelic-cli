package credentials

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
)

var defaultProfile *Profile

// WithCredentials loads and returns the CLI credentials.
func WithCredentials(f func(c *Credentials)) {
	WithCredentialsFrom(config.DefaultConfigDirectory, f)
}

// WithCredentialsFrom loads and returns the CLI credentials from a specified location.
func WithCredentialsFrom(configDir string, f func(c *Credentials)) {
	c, err := LoadCredentials(configDir)
	if err != nil {
		log.Fatal(err)
	}

	f(c)
}

// DefaultProfile retrieves the current default profile.
func DefaultProfile() *Profile {
	if defaultProfile == nil {
		WithCredentials(func(c *Credentials) {
			defaultProfile = c.Default()
		})
	}

	return defaultProfile
}

// SetDefaultProfile allows mocking of the default profile for testing purposes.
func SetDefaultProfile(p Profile) {
	defaultProfile = &p
}
