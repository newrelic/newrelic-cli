package servicelevel

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicelevel"
)

var (
	createEntityGUID  string
	createName        string
	createDescription string
	createAccountID   int
	createFacets      []string
	createValidFrom   string
	createValidWhere  string
	createGoodFrom    string
	createGoodWhere   string
	createBadFrom     string
	createBadWhere    string
)

var CmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new Service Level Indicator",
	Long:  `Create a new Service Level Indicator via the New Relic API.`,
	Example: `# Create a Service Level Indicator
newrelic servicelevel create --entityGuid <guid> --name <name> --accountId <accountId> --facets <facets>`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		eventsInput := servicelevel.ServiceLevelEventsCreateInput{
			AccountID: createAccountID,
			Facets:    createFacets,
		}

		if createValidFrom != "" {
			eventsInput.ValidEvents = &servicelevel.ServiceLevelEventsQueryCreateInput{
				From:  servicelevel.NRQL(createValidFrom),
				Where: servicelevel.NRQL(createValidWhere),
			}
		}

		if createGoodFrom != "" {
			eventsInput.GoodEvents = &servicelevel.ServiceLevelEventsQueryCreateInput{
				From:  servicelevel.NRQL(createGoodFrom),
				Where: servicelevel.NRQL(createGoodWhere),
			}
		}

		if createBadFrom != "" {
			eventsInput.BadEvents = &servicelevel.ServiceLevelEventsQueryCreateInput{
				From:  servicelevel.NRQL(createBadFrom),
				Where: servicelevel.NRQL(createBadWhere),
			}
		}

		input := servicelevel.ServiceLevelIndicatorCreateInput{
			Name:        createName,
			Description: createDescription,
			Events:      eventsInput,
		}

		result, err := client.NRClient.ServiceLevel.ServiceLevelCreateWithContext(
			utils.SignalCtx,
			common.EntityGUID(createEntityGUID),
			input,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdCreate)

	CmdCreate.Flags().StringVar(&createEntityGUID, "entityGuid", "", "The entity GUID to attach the SLI to")
	utils.LogIfError(CmdCreate.MarkFlagRequired("entityGuid"))

	CmdCreate.Flags().StringVar(&createName, "name", "", "The name of the SLI")
	utils.LogIfError(CmdCreate.MarkFlagRequired("name"))

	CmdCreate.Flags().StringVar(&createDescription, "description", "", "The description of the SLI")

	CmdCreate.Flags().IntVar(&createAccountID, "accountId", 0, "The New Relic account ID where the events are fetched from")
	utils.LogIfError(CmdCreate.MarkFlagRequired("accountId"))

	CmdCreate.Flags().StringSliceVar(&createFacets, "facets", []string{}, "The list of NRQL FACET attributes to group results by")
	utils.LogIfError(CmdCreate.MarkFlagRequired("facets"))

	CmdCreate.Flags().StringVar(&createValidFrom, "validEventsFrom", "", "The NRDB event to fetch valid events data from (NRQL FROM clause)")
	CmdCreate.Flags().StringVar(&createValidWhere, "validEventsWhere", "", "The NRQL condition to filter valid events")

	CmdCreate.Flags().StringVar(&createGoodFrom, "goodEventsFrom", "", "The NRDB event to fetch good events data from (NRQL FROM clause)")
	CmdCreate.Flags().StringVar(&createGoodWhere, "goodEventsWhere", "", "The NRQL condition to filter good events")

	CmdCreate.Flags().StringVar(&createBadFrom, "badEventsFrom", "", "The NRDB event to fetch bad events data from (NRQL FROM clause)")
	CmdCreate.Flags().StringVar(&createBadWhere, "badEventsWhere", "", "The NRQL condition to filter bad events")
}
