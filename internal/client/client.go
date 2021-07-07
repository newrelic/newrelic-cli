package client

import (
	"errors"
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	NRClient    *newrelic.NewRelic
	serviceName = "newrelic-cli"
	version     = "dev"
)

// NewClient initializes the New Relic client.
func NewClient(profileName string) (*newrelic.NewRelic, error) {
	apiKey := configuration.GetProfileString(profileName, configuration.APIKey)
	insightsInsertKey := configuration.GetProfileString(profileName, configuration.InsightsInsertKey)

	if apiKey == "" && insightsInsertKey == "" {
		return nil, errors.New("a User API key or Ingest API key is required, set a default profile or use the NEW_RELIC_API_KEY or NEW_RELIC_INSIGHTS_INSERT_KEY environment variables")
	}

	region := configuration.GetProfileString(profileName, configuration.Region)
	logLevel := configuration.GetConfigString(configuration.LogLevel)
	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	cfgOpts := []newrelic.ConfigOption{
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigInsightsInsertKey(insightsInsertKey),
		newrelic.ConfigLogLevel(logLevel),
		newrelic.ConfigRegion(region),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	}

	nrClient, err := newrelic.New(cfgOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}

func RequireClient(cmd *cobra.Command, args []string) {
	if NRClient == nil {
		log.Fatalf("could not initialize New Relic client, make sure your profile is configured with `newrelic profile configure`")
	}
}
