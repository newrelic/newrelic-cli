package main

import (
	"github.com/newrelic/newrelic-client-go/newrelic"
	log "github.com/sirupsen/logrus"

	// Commands
	root "github.com/newrelic/newrelic-cli/internal/cmd"
	"github.com/newrelic/newrelic-cli/internal/entities"

	"github.com/newrelic/newrelic-cli/pkg/apm"
	"github.com/newrelic/newrelic-cli/pkg/config"
	"github.com/newrelic/newrelic-cli/pkg/credentials"
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
	root.Command.AddCommand(entities.Command)
	root.Command.AddCommand(credentials.Command)
	root.Command.AddCommand(apm.Command)
	root.Command.AddCommand(config.Command)
}

func main() {
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}

	if err := loadCredentials(); err != nil {
		log.Fatal(err)
	}

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

	if err := root.Execute(AppName, Version); err != nil {
		log.Fatal(err)
	}
}
