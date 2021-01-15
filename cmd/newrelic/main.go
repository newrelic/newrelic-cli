package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/agent"
	"github.com/newrelic/newrelic-cli/internal/apiaccess"
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/config"
	diagnose "github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/edge"
	"github.com/newrelic/newrelic-cli/internal/entities"
	"github.com/newrelic/newrelic-cli/internal/events"
	"github.com/newrelic/newrelic-cli/internal/install"
	"github.com/newrelic/newrelic-cli/internal/nerdgraph"
	"github.com/newrelic/newrelic-cli/internal/nerdstorage"
	"github.com/newrelic/newrelic-cli/internal/nrql"
	"github.com/newrelic/newrelic-cli/internal/profiles"
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
	Command.AddCommand(diagnose.Command)
	Command.AddCommand(edge.Command)
	Command.AddCommand(events.Command)
	Command.AddCommand(entities.Command)
	Command.AddCommand(nerdgraph.Command)
	Command.AddCommand(nerdstorage.Command)
	Command.AddCommand(nrql.Command)
	Command.AddCommand(profiles.Command)
	Command.AddCommand(reporting.Command)
	Command.AddCommand(workload.Command)
	Command.AddCommand(agent.Command)
	Command.AddCommand(install.Command)
	Command.AddCommand(install.TestCommand)
	Command.AddCommand(apiaccess.Command)

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
	v, err := config.GetConfigValue(config.PrereleaseFeatures)
	if err != nil {
		log.Fatal(err)
	}

	if !config.Ternary(v.(string)).Bool() {
		return
	}

	log.Debug("Pre-release mode active")

	for _, cmd := range c.Commands() {
		if cmd.Hidden {
			log.Debugf("Activating pre-release subcommand: %s", cmd.Name())
			cmd.Hidden = false
		}
	}
}
