package client

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/region"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
)

var (
	serviceName = "newrelic-cli"
	version     = "dev"
)

// CreateNRClient initializes the New Relic client.
func CreateNRClient(cfg *config.Config, creds *credentials.Credentials) (*newrelic.NewRelic, error) {
	var (
		err         error
		apiKey      string
		regionValue region.Name
	)

	// Create the New Relic Client
	defProfile := creds.Default()

	defProfile = applyOverrides(defProfile)

	if defProfile != nil {
		apiKey = defProfile.APIKey
		regionValue = defProfile.Region
	}

	if apiKey == "" {
		return nil, errors.New("an API key is required, set a default profile or use the NEW_RELIC_API_KEY environment variable")
	}

	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(regionValue),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}

// applyOverrides reads Profile info out of the Environment to override config
func applyOverrides(p *credentials.Profile) *credentials.Profile {
	envAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	envRegion := os.Getenv("NEW_RELIC_REGION")

	if envAPIKey == "" && envRegion == "" {
		return p
	}

	out := credentials.Profile{}
	if p != nil {
		out = *p
	}

	if envAPIKey != "" {
		out.APIKey = envAPIKey
	}

	if envRegion != "" {
		var err error
		out.Region, err = region.Parse(envRegion)

		if err != nil {
			switch err.(type) {
			case region.UnknownError:
				log.Errorf("error parsing NEW_RELIC_REGION: %s", err)
				// Ignore the override if they have a default on the profile
				if p.Region != "" {
					var e2 error
					out.Region, e2 = region.Parse(p.Region.String())
					if e2 != nil {
						log.Errorf("error parsing default profile: %s", e2)
						out.Region = region.Default
					}
				} else {
					out.Region = region.Default
				}
				log.Errorf("using region %s", out.Region.String())
			case region.UnknownUsingDefaultError:
				log.Error(err)
			default:
				log.Fatalf("unknown error: %v", err)
			}
		}
	}

	return &out
}
