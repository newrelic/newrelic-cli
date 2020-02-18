package main

import (
	"github.com/newrelic/newrelic-client-go/newrelic"
	log "github.com/sirupsen/logrus"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/entities"
)

var (
	// AppName for this CMD
	AppName = "newrelic-dev"
	// Version of the CLI
	Version = "dev"

	globalConfig *config.Config
	creds        *credentials.Credentials

	// Client is an instance of the New Relic client.
	nrClient *newrelic.NewRelic
)

func init() {
	// Bind imported sub-commands
	Command.AddCommand(entities.Command)
	Command.AddCommand(credentials.Command)
	Command.AddCommand(apm.Command)
	Command.AddCommand(config.Command)
}

func main() {
	// TODO Here too we should probably return the client rather than reaching
	// into the global.
	if err := createNRClient(); err != nil {
		log.Fatal(err)
	}

	// Configure commands that need it
	entities.SetClient(nrClient)
	apm.SetClient(nrClient)

	credentials.SetConfig(globalConfig)
	config.SetConfig(globalConfig)
	credentials.SetCredentials(creds)

	if err := Execute(AppName, Version); err != nil {
		log.Fatal(err)
	}
}
