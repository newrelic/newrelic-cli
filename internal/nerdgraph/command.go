package nerdgraph

import (
	"github.com/spf13/cobra"
)

// Command represents the nerdgraph command.
var Command = &cobra.Command{
	Use:   "nerdgraph",
	Short: "Subcommand for executing raw GraphQL requests to the NerdGraph API.",
}
