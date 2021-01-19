package events

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/utils"
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
		accountID = config.FatalIfAccountIDNotPresent()
		config.FatalIfActiveProfileFieldStringNotPresent(config.InsightsInsertKey)
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
	utils.LogIfError(cmdPost.MarkFlagRequired("event"))
}
