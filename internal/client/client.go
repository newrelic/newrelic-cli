package client

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-client-go/newrelic"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
)

var (
	serviceName = "newrelic-cli"
	version     = "dev"
)

// CreateNRClient initializes the New Relic client.
func CreateNRClient(cfg *config.Config, creds *credentials.Credentials) (*newrelic.NewRelic, *credentials.Profile, error) {
	var (
		err               error
		apiKey            string
		insightsInsertKey string
		regionValue       string
	)

	// Create the New Relic Client
	defProfile := creds.Default()

	defProfile = applyOverrides(defProfile)

	if defProfile != nil {
		apiKey = defProfile.APIKey
		insightsInsertKey = defProfile.InsightsInsertKey
		regionValue = defProfile.Region
	}

	if apiKey == "" {
		return nil, nil, errors.New("an API key is required, set a default profile or use the NEW_RELIC_API_KEY environment variable")
	}

	userAgent := fmt.Sprintf("newrelic-cli/%s (https://github.com/newrelic/newrelic-cli)", version)

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(apiKey),
		newrelic.ConfigInsightsInsertKey(insightsInsertKey),
		newrelic.ConfigLogLevel(cfg.LogLevel),
		newrelic.ConfigRegion(regionValue),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, defProfile, nil
}

// applyOverrides reads Profile info out of the Environment to override config
func applyOverrides(p *credentials.Profile) *credentials.Profile {
	envAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	envInsightsInsertKey := os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")
	envRegion := os.Getenv("NEW_RELIC_REGION")
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")

	if envAPIKey == "" && envRegion == "" && envInsightsInsertKey == "" && envAccountID == "" {
		return p
	}

	out := credentials.Profile{}
	if p != nil {
		out = *p
	}

	if envAPIKey != "" {
		out.APIKey = envAPIKey
	}

	if envInsightsInsertKey != "" {
		out.InsightsInsertKey = envInsightsInsertKey
	}

	if envRegion != "" {
		out.Region = strings.ToUpper(envRegion)
	}

	if envAccountID != "" {
		accountID, err := strconv.Atoi(envAccountID)
		if err != nil {
			log.Warnf("Invalid account ID: %s", envAccountID)
			return &out
		}

		out.AccountID = accountID
	}

	return &out
}
