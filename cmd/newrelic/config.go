package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/pkg/config"
)

// loadConfig loads the configuration
func loadConfig() error {
	var err error

	globalConfig, err = config.Load(configFile, logLevel)
	if err != nil {
		return err
	}

	log.Tracef("config: %+v\n", globalConfig)

	return nil
}
