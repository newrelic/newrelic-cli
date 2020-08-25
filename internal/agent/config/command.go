package apm

import (
	"github.com/spf13/cobra"
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "agent",
	Short: "Utilities for New Relic Agents",
	Long:  `Utilities for New Relic Agents`,
}

var cmdConfig = &cobra.Command{
	Use:     "config",
	Short:   "Configuration Utilities/Helpers for New Relic Agents",
	Long:    `Utilities and Helpers to aid in the configuration of New Relic Agents.`,
	Example: "newrelic agent config obfuscate --key <key> --value <value>",
}

func init() {

	Command.AddCommand(cmdConfig)

}
