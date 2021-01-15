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
	accountID int
	event     string
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
	Run: func(cmd *cobra.Command, args []string) {
		insightsInsertKey := config.GetActiveProfileValueString(config.InsightsInsertKey)
		if insightsInsertKey == "" {
			log.Fatal("an Insights insert key is required, set one in your default profile or use the NEW_RELIC_INSIGHTS_INSERT_KEY environment variable")
		}

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
	cmdPost.Flags().IntVarP(&accountID, "accountId", "a", 0, "the account ID to create the custom event in")
	cmdPost.Flags().StringVarP(&event, "event", "e", "{}", "a JSON-formatted event payload to post")
	utils.LogIfError(cmdPost.MarkFlagRequired("accountId"))
	utils.LogIfError(cmdPost.MarkFlagRequired("event"))
}
