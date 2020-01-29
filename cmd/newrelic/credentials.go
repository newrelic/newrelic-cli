package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/pkg/config"
	"github.com/newrelic/newrelic-cli/pkg/profile"
)

// loadCredentials loads the list of profiles
func loadCredentials() error {
	var err error

	if globalConfig == nil {
		if err = loadConfig(); err != nil {
			return nil
		}
	}

	// Load profiles
	creds, err = profile.Load(config.DefaultConfigDirectory, config.DefaultEnvPrefix)
	if err != nil {
		// TODO: Don't die here, we need to be able to run the profiles command to add one!
		return err
	}

	log.Tracef("profiles: %+v\n", creds)

	return nil
}
