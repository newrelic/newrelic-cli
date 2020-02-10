package main

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/newrelic"
)

// createNRClient initializes the Client global
func createNRClient() error {
	var (
		err            error
		personalAPIKey string
		adminAPIKey    string
		region         string
	)

	if creds == nil {
		if err = loadCredentials(); err != nil {
			return err
		}
	}

	// Create the New Relic Client
	defProfile := creds.Default()
	if defProfile != nil {
		adminAPIKey = defProfile.AdminAPIKey
		personalAPIKey = defProfile.PersonalAPIKey
		region = defProfile.Region
	} else {
		return fmt.Errorf("invalid profile name: '%s'", creds.DefaultProfile)
	}

	nrClient, err = newrelic.New(
		newrelic.ConfigAPIKey(adminAPIKey),
		newrelic.ConfigPersonalAPIKey(personalAPIKey),
		newrelic.ConfigLogLevel(globalConfig.LogLevel),
		newrelic.ConfigRegion(region),
	)
	if err != nil {
		return fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nil
}
