package edge

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/pkg/edge"
)

var (
	accountID            int
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
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if accountID, err = config.RequireAccountID(); err != nil {
			log.Fatal(err)
		}

		if _, err = config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the New Relic Edge trace observers for an account.",
	Long: `List the New Relic trace observers for an account

The list command retrieves the trace observers for the given account ID.
`,
	Example: `newrelic edge trace-observer list --accountId 12345678`,
	Run: func(cmd *cobra.Command, args []string) {
		traceObservers, err := client.Client.Edge.ListTraceObservers(accountID)
		if err != nil {
			log.Fatal(err)
		}

		if err := output.Print(traceObservers); err != nil {
			log.Fatal(err)
		}
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
	Run: func(cmd *cobra.Command, args []string) {
		if ok := isValidProviderRegion(providerRegion); !ok {
			log.Fatalf("%s is not a valid provider region, valid values are %s", providerRegion, validProviderRegions)
		}

		traceObserver, err := client.Client.Edge.CreateTraceObserver(accountID, name, edge.EdgeProviderRegion(providerRegion))
		if err != nil {
			log.Fatal(err)
		}

		if err = output.Print(traceObserver); err != nil {
			log.Fatal(err)
		}

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
	Run: func(cmd *cobra.Command, args []string) {
		_, err := client.Client.Edge.DeleteTraceObserver(accountID, id)
		if err != nil {
			log.Fatal(err)
		}

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
	if err := cmdCreate.MarkFlagRequired("name"); err != nil {
		log.Error(err)
	}

	if err := cmdCreate.MarkFlagRequired("providerRegion"); err != nil {
		log.Error(err)
	}

	// Delete
	cmdTraceObserver.AddCommand(cmdDelete)
	cmdDelete.Flags().IntVarP(&id, "id", "i", 0, "the ID of the trace observer to delete")
	if err := cmdDelete.MarkFlagRequired("id"); err != nil {
		log.Error(err)
	}
}
