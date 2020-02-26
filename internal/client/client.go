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
		adminAPIKey    string
		region         string
	)

	// Create the New Relic Client
	defProfile := creds.Default()

	defProfile = applyOverrides(defProfile)

	if defProfile != nil {
		adminAPIKey = defProfile.AdminAPIKey
		personalAPIKey = defProfile.PersonalAPIKey
		region = defProfile.Region
	} else {
		return nil, fmt.Errorf("invalid profile name: '%s'", creds.DefaultProfile)
	}

	nrClient, err := newrelic.New(
		newrelic.ConfigAPIKey(adminAPIKey),
		newrelic.ConfigPersonalAPIKey(personalAPIKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}

func applyOverrides(p *credentials.Profile) *credentials.Profile {
	envAPIKey := os.Getenv("NEWRELIC_API_KEY")
	envAdminAPIKey := os.Getenv("NEWRELIC_ADMIN_API_KEY")
	envRegion := os.Getenv("NEWRELIC_REGION")

	if envAPIKey == "" && envAdminAPIKey == "" && envRegion == "" {
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

	if envAdminAPIKey != "" {
		out.AdminAPIKey = envAdminAPIKey
	}

	if envRegion != "" {
		out.Region = envRegion
	}

	return &out
}
