package client

import (
	"fmt"
	"os"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

// CreateNRClient initializes the New Relic client.
func CreateNRClient(cfg *config.Config, creds *credentials.Credentials) (*newrelic.NewRelic, error) {
	var (
		err            error
		personalAPIKey string
		region         string
	)

	// Create the New Relic Client
	defProfile := creds.Default()

	defProfile = applyOverrides(defProfile)

	if defProfile != nil {
		personalAPIKey = defProfile.PersonalAPIKey
		region = defProfile.Region
	} else {
		return nil, fmt.Errorf("invalid profile name: '%s'", creds.DefaultProfile)
	}

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(personalAPIKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}

// applyOverrides reads Profile info out of the Environment to override config
func applyOverrides(p *credentials.Profile) *credentials.Profile {
	envAPIKey := os.Getenv("NEWRELIC_PERSONAL_API_KEY")
	envRegion := os.Getenv("NEWRELIC_REGION")

	if envAPIKey == "" && envRegion == "" {
		return p
	}

	var out credentials.Profile

	if p == nil {
		out = credentials.Profile{}
	} else {
		out = *p
	}

	if envAPIKey != "" {
		out.PersonalAPIKey = envAPIKey
	}

	if envRegion != "" {
		out.Region = envRegion
	}

	return &out
}
