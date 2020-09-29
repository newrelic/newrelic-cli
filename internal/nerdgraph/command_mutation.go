package nerdgraph

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

var cmdMutation = &cobra.Command{
	Use:     "mutation",
	Short:   "here is a short description",
	Long:    "here is a longer description with more detail",
	Example: "newrelic nerdgraph mutation --help",
}

var accountID int
var alertsPolicyCreateInput string

var cmdAlertsPolicyCreate = &cobra.Command{
	Use:     "alertsPolicyCreate",
	Short:   "here is a short description",
	Long:    "here is a longer description with more detail",
	Example: "newrelic nerdgraph mutation alertsPolicyCreate --input='{\"name\": \"foo\",\"incidentPreference\": \"PER_CONDITION\"}' --accountId=$NEW_RELIC_ACCOUNT_ID\n",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var input alerts.AlertsPolicyInput
			err := json.Unmarshal([]byte(alertsPolicyCreateInput), &input)
			utils.LogIfFatal(err)

			resp, err := nrClient.Alerts.CreatePolicyMutation(accountID, input)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(resp))
		})
	},
}

func init() {
	Command.AddCommand(cmdMutation)

	cmdMutation.AddCommand(cmdAlertsPolicyCreate)

	cmdAlertsPolicyCreate.Flags().IntVar(&accountID, "accountId", 0, "describe the flag here")
	utils.LogIfError(cmdAlertsPolicyCreate.MarkFlagRequired("accountId"))

	cmdAlertsPolicyCreate.Flags().StringVar(&alertsPolicyCreateInput, "input", "", "describe the flag here")
}
