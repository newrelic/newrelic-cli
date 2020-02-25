package client

import (
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
	log "github.com/sirupsen/logrus"
)

// WithClient returns a New Relic client.
func WithClient(f func(c *newrelic.NewRelic)) {
	WithClientFrom(config.DefaultConfigDirectory, f)
}

// WithClientFrom returns a New Relic client, initialized from configuration in the specified location.
func WithClientFrom(configDir string, f func(c *newrelic.NewRelic)) {
	config.WithConfigFrom(configDir, func(cfg *config.Config) {
		credentials.WithCredentialsFrom(configDir, func(creds *credentials.Credentials) {
			nrClient, err := CreateNRClient(cfg, creds)
			if err != nil {
				log.Fatal("cannot initialize client")
			}

			f(nrClient)
		})
	})
}
