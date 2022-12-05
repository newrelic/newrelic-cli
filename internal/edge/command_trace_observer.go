package edge

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/edge"
)

var (
	id                   int
	name                 string
	providerRegion       string
	validProviderRegions = []string{
		string(edge.EdgeProviderRegionTypes.AWS_US_EAST_1),
		string(edge.EdgeProviderRegionTypes.AWS_US_EAST_2),
	}
)

var cmdTraceObserver = &cobra.Command{
	Use:   "trace-observer",
	Short: "Interact with New Relic Edge trace observers.",
	Long: `Interact with New Relic Edge trace observers
	
	A trace observer is a configuration that enables infinite tracing for an account.
	Once enabled, infinite tracing observes 100% of your application traces, then
	provides visualization for the most actionable data so you can investigate and
	solve issues faster.`,
	Example: "newrelic edge trace-observer list --accountId <accountID>",
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the New Relic Edge trace observers for an account.",
	Long: `List the New Relic trace observers for an account

The list command retrieves the trace observers for the given account ID.
`,
	Example: `newrelic edge trace-observer list --accountId 12345678`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()
		traceObservers, err := client.NRClient.Edge.ListTraceObserversWithContext(utils.SignalCtx, accountID)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(traceObservers))
	},
}

func isValidProviderRegion(providerRegion string) bool {
	for _, r := range validProviderRegions {
		if r == providerRegion {
			return true
		}
	}

	return false
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic Edge trace observer.",
	Long: `Create a New Relic Edge trace observer

The create command requires an account ID, observer name, and provider region.
Valid provider regions are AWS_US_EAST_1 and AWS_US_EAST_2.
`,
	Example: `newrelic edge trace-observer create --name 'My Observer' --accountId 12345678 --providerRegion AWS_US_EAST_1`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()
		if ok := isValidProviderRegion(providerRegion); !ok {
			log.Fatalf("%s is not a valid provider region, valid values are %s", providerRegion, validProviderRegions)
		}

		traceObserver, err := client.NRClient.Edge.CreateTraceObserver(accountID, name, edge.EdgeProviderRegion(providerRegion))
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(traceObserver))
		log.Info("success")
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a New Relic Edge trace observer.",
	Long: `Delete a New Relic Edge trace observer.

The delete command accepts a trace observer's ID.
`,
	Example: `newrelic edge trace-observer delete --accountId 12345678 --id 1234`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()
		_, err := client.NRClient.Edge.DeleteTraceObserver(accountID, id)
		utils.LogIfFatal(err)

		log.Info("success")
	},
}

func init() {
	// Root sub-command
	Command.AddCommand(cmdTraceObserver)

	// List
	cmdTraceObserver.AddCommand(cmdList)

	// Create
	cmdTraceObserver.AddCommand(cmdCreate)
	cmdCreate.Flags().StringVarP(&name, "name", "n", "", "the name of the trace observer")
	cmdCreate.Flags().StringVarP(&providerRegion, "providerRegion", "r", "", "the provider region in which to create the trace observer")
	utils.LogIfError(cmdCreate.MarkFlagRequired("name"))
	utils.LogIfError(cmdCreate.MarkFlagRequired("providerRegion"))

	// Delete
	cmdTraceObserver.AddCommand(cmdDelete)
	cmdDelete.Flags().IntVarP(&id, "id", "i", 0, "the ID of the trace observer to delete")
	utils.LogIfError(cmdDelete.MarkFlagRequired("id"))
}
