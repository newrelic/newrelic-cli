package events

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
)

var (
	event     string
	accountID int
)

var cmdPost = &cobra.Command{
	Use:   "post",
	Short: "Post a custom event to New Relic",
	Long: `Post a custom event to New Relic

The post command accepts an account ID and JSON-formatted payload representing a
custom event to be posted to New Relic. These events once posted can be queried
using NRQL via the CLI or New Relic One UI.
The accepted payload requires the use of an ` + "`eventType`" + `field that
represents the custom event's type.
`,
	Example: `newrelic events post --accountId 12345 --event '{ "eventType": "Payment", "amount": 123.45 }'`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if accountID, err = config.RequireAccountID(); err != nil {
			log.Fatal(err)
		}

		if _, err = config.RequireInsightsInsertKey(); err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var e map[string]interface{}

		if err := json.Unmarshal([]byte(event), &e); err != nil {
			log.Fatal(err)
		}

		if err := client.Client.Events.CreateEvent(accountID, event); err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

func init() {
	Command.AddCommand(cmdPost)
	cmdPost.Flags().StringVarP(&event, "event", "e", "{}", "a JSON-formatted event payload to post")
	if err := cmdPost.MarkFlagRequired("event"); err != nil {
		log.Error(err)
	}
}
