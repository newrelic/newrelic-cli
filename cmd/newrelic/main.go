package main

import (
	"github.com/newrelic/newrelic-client-go/newrelic"
	log "github.com/sirupsen/logrus"

	// Commands
	root "github.com/newrelic/newrelic-cli/internal/cmd"
	"github.com/newrelic/newrelic-cli/internal/entities"

	"github.com/newrelic/newrelic-cli/pkg/config"
	"github.com/newrelic/newrelic-cli/pkg/profile"
)

var (
	// AppName for this CMD
	AppName = "newrelic-dev"
	// Version of the CLI
	Version = "dev"

	globalConfig *config.Config
	creds        *profile.Credentials

	// Client is an instance of the New Relic client.
	nrClient *newrelic.NewRelic

	configFile string
	logLevel   string
)

func init() {
	if err := loadConfig(); err != nil {
		log.Fatal(err)
	}

	if err := loadCredentials(); err != nil {
		log.Fatal(err)
	}

	if err := createNRClient(); err != nil {
		log.Fatal(err)
	}

	// We want to track these at the global level, but need them here...
	root.Command.PersistentFlags().StringVar(&configFile, "config", "", "config file (default: $HOME/.newrelic/config.json)")
	root.Command.PersistentFlags().StringVar(&logLevel, "log-level", "", "log level [Panic,Fatal,Error,Warn,Info,Debug,Trace]")

	// Bind imported sub-commands
	root.Command.AddCommand(entities.Command)
	root.Command.AddCommand(profile.Command)
}

func main() {
	// Configure commands that need it
	entities.SetClient(nrClient)

	profile.SetConfig(globalConfig)
	profile.SetCredentials(creds)

	if err := root.Execute(AppName, Version); err != nil {
		log.Fatal(err)
	}
}
