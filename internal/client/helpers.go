package client

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

// WithClient returns a New Relic client.
func WithClient(f func(c *newrelic.NewRelic)) {
	WithClientFrom(config.DefaultConfigDirectory, f)
}

// WithClientFrom returns a New Relic client, initialized from configuration in the specified location.
func WithClientFrom(configDir string, f func(c *newrelic.NewRelic)) {
	config.WithConfigFrom(configDir, func(cfg *config.Config) {
		credentials.WithCredentialsFrom(configDir, func(creds *credentials.Credentials) {
			nrClient, _, err := CreateNRClient(cfg, creds)
			if err != nil {
				log.Fatal(err)
			}

			f(nrClient)
		})
	})
}

// WithClientAndProfile returns a New Relic client and the profile used to initialize it,
// after environment oveerrides have been applied.
func WithClientAndProfile(f func(c *newrelic.NewRelic, p *credentials.Profile)) {
	WithClientAndProfileFrom(config.DefaultConfigDirectory, f)
}

// WithClientAndProfileFrom returns a New Relic client and default profile used to initialize it,
// after environment oveerrides have been applied.
func WithClientAndProfileFrom(configDir string, f func(c *newrelic.NewRelic, p *credentials.Profile)) {
	config.WithConfigFrom(configDir, func(cfg *config.Config) {
		credentials.WithCredentialsFrom(configDir, func(creds *credentials.Credentials) {
			nrClient, defaultProfile, err := CreateNRClient(cfg, creds)
			if err != nil {
				log.Fatal(err)
			}

			f(nrClient, defaultProfile)
		})
	})
}
