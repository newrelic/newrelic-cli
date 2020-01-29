package main

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/newrelic"
)

// createNRClient initializes the Client global
func createNRClient() error {
	var (
		err    error
		apiKey string
		region string
	)

	if profiles == nil {
		if err = loadProfiles(); err != nil {
			return err
		}
	}

	// Create the New Relic Client
	defProfile := profiles.Default()
	if defProfile != nil {
		apiKey = defProfile.PersonalAPIKey
		region = defProfile.Region
	} else {
		return fmt.Errorf("invalid profile name: '%s'", profiles.DefaultName)
	}

	nrClient, err = newrelic.New(newrelic.ConfigPersonalAPIKey(apiKey), newrelic.ConfigLogLevel(globalConfig.LogLevel), newrelic.ConfigRegion(region))
	if err != nil {
		return fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nil
}
