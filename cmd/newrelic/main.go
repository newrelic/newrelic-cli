package main

import (
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/entities"
	"github.com/newrelic/newrelic-cli/internal/nerdgraph"
	"github.com/newrelic/newrelic-cli/internal/nerdstorage"
	"github.com/newrelic/newrelic-cli/internal/workload"
)

var (
	appName = "newrelic-dev"
	version = "dev"
)

func init() {
	// Bind imported sub-commands
	Command.AddCommand(entities.Command)
	Command.AddCommand(credentials.Command)
	Command.AddCommand(apm.Command)
	Command.AddCommand(config.Command)
	Command.AddCommand(nerdgraph.Command)
	Command.AddCommand(nerdstorage.Command)
	Command.AddCommand(workload.Command)
}

func main() {
	if err := Execute(); err != nil {
		if err != flag.ErrHelp {
			log.Fatal(err)
		}
	}
}
