package nerdstorage

import (
	log "github.com/sirupsen/logrus"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if accountID, err = config.RequireAccountID(); err != nil {
			log.Fatal(err)
		}

		if _, err = config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
}
