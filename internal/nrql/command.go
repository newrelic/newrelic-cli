package nrql

import (
	"github.com/spf13/cobra"
)

// Command represents the nerdgraph command.
var Command = &cobra.Command{
	Use:   "nrql",
	Short: "Commands for interacting with the New Relic Database",
}
