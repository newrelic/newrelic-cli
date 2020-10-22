package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/agent"
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/edge"
	"github.com/newrelic/newrelic-cli/internal/entities"
	"github.com/newrelic/newrelic-cli/internal/events"
	"github.com/newrelic/newrelic-cli/internal/install"
	"github.com/newrelic/newrelic-cli/internal/nerdgraph"
	"github.com/newrelic/newrelic-cli/internal/nerdstorage"
	"github.com/newrelic/newrelic-cli/internal/nrql"
	"github.com/newrelic/newrelic-cli/internal/reporting"
	"github.com/newrelic/newrelic-cli/internal/workload"
)

var (
	appName = "newrelic-dev"
	version = "dev"
)

func init() {
	// Bind imported sub-commands
	Command.AddCommand(apm.Command)
	Command.AddCommand(config.Command)
	Command.AddCommand(credentials.Command)
	Command.AddCommand(edge.Command)
	Command.AddCommand(events.Command)
	Command.AddCommand(entities.Command)
	Command.AddCommand(nerdgraph.Command)
	Command.AddCommand(nerdstorage.Command)
	Command.AddCommand(nrql.Command)
	Command.AddCommand(reporting.Command)
	Command.AddCommand(workload.Command)
	Command.AddCommand(agent.Command)
	Command.AddCommand(install.Command)

	CheckPrereleaseMode(Command)
}

func main() {
	if err := Execute(); err != nil {
		if err != flag.ErrHelp {
			log.Fatal(err)
		}
	}
}

// CheckPrereleaseMode unhides subcommands marked as hidden when the pre-release
// flag is active.
func CheckPrereleaseMode(c *cobra.Command) {
	config.WithConfig(func(cfg *config.Config) {
		if !cfg.PreReleaseFeatures.Bool() {
			return
		}

		log.Debug("Pre-release mode active")

		for _, cmd := range c.Commands() {
			if cmd.Hidden {
				log.Debugf("Activating pre-release subcommand: %s", cmd.Name())
				cmd.Hidden = false
			}
		}
	})
}
