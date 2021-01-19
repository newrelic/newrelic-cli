package nerdstorage

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
)

var (
	accountID  int
	entityGUID string
	packageID  string
	collection string
	documentID string
	document   string
	scope      string
)

// Command represents the nerdstorage command.
var Command = &cobra.Command{
	Use:   "nerdstorage",
	Short: "Read, write, and delete NerdStorage documents and collections.",
	PreRun: func(cmd *cobra.Command, args []string) {
		accountID = config.FatalIfAccountIDNotPresent()
		config.FatalIfActiveProfileFieldStringNotPresent(config.APIKey)
	},
}
