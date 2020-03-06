package apm

import (
	"github.com/spf13/cobra"
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "apm",
	Short: "Subcommands to interact with New Relic APM",
}
