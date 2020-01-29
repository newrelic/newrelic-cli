package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/pkg/config"
	"github.com/newrelic/newrelic-cli/pkg/profile"
)

// loadProfiles loads the list of profiles
func loadProfiles() error {
	var err error

	if globalConfig == nil {
		if err = loadConfig(); err != nil {
			return nil
		}
	}

	// Load profiles
	profiles, err = profile.Load(config.DefaultConfigDirectory, config.DefaultEnvPrefix)
	if err != nil {
		// TODO: Don't die here, we need to be able to run the profiles command to add one!
		return err
	}

	log.Tracef("profiles: %+v\n", profiles)

	return nil
}
