//go:generate goversioninfo -o=resource_windows.syso
package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	// Commands
	"github.com/newrelic/newrelic-cli/internal/agent"
	"github.com/newrelic/newrelic-cli/internal/apiaccess"
	"github.com/newrelic/newrelic-cli/internal/apm"
	"github.com/newrelic/newrelic-cli/internal/cli"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	configCmd "github.com/newrelic/newrelic-cli/internal/config/command"
	"github.com/newrelic/newrelic-cli/internal/decode"
	diagnose "github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/edge"
	"github.com/newrelic/newrelic-cli/internal/entities"
	"github.com/newrelic/newrelic-cli/internal/events"
	"github.com/newrelic/newrelic-cli/internal/install"
	"github.com/newrelic/newrelic-cli/internal/nerdgraph"
	"github.com/newrelic/newrelic-cli/internal/nerdstorage"
	"github.com/newrelic/newrelic-cli/internal/nrql"
	"github.com/newrelic/newrelic-cli/internal/profile"
	"github.com/newrelic/newrelic-cli/internal/reporting"
	"github.com/newrelic/newrelic-cli/internal/synthetics"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-cli/internal/workload"
)

var (
	appName = "newrelic-dev"
)

func init() {
	// Bind imported sub-commands
	Command.AddCommand(agent.Command)
	Command.AddCommand(apiaccess.Command)
	Command.AddCommand(synthetics.Command)
	Command.AddCommand(apm.Command)
	Command.AddCommand(configCmd.Command)
	Command.AddCommand(decode.Command)
	Command.AddCommand(diagnose.Command)
	Command.AddCommand(edge.Command)
	Command.AddCommand(entities.Command)
	Command.AddCommand(events.Command)
	Command.AddCommand(install.Command)
	Command.AddCommand(nerdgraph.Command)
	Command.AddCommand(nerdstorage.Command)
	Command.AddCommand(nrql.Command)
	Command.AddCommand(profile.Command)
	Command.AddCommand(reporting.Command)
	Command.AddCommand(utils.Command)
	Command.AddCommand(workload.Command)

	CheckPrereleaseMode(Command)

	os.Setenv("NEW_RELIC_CLI_VERSION", cli.Version())
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
	if configAPI.GetConfigTernary(config.PreReleaseFeatures).Bool() {
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
