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
	updateGUID        string
	updateID          string
	updateName        string
	updateDescription string
	updateFacets      []string
	updateValidFrom   string
	updateValidWhere  string
	updateGoodFrom    string
	updateGoodWhere   string
	updateBadFrom     string
	updateBadWhere    string
)

var CmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update an existing Service Level Indicator",
	Long:  `Update an existing Service Level Indicator via the New Relic API.`,
	Example: `# Update a Service Level Indicator
newrelic servicelevel update --guid <guid> --id <id> --name <name>`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		eventsInput := servicelevel.ServiceLevelEventsUpdateInput{
			Facets: updateFacets,
		}

		if updateValidFrom != "" {
			eventsInput.ValidEvents = &servicelevel.ServiceLevelEventsQueryUpdateInput{
				From:  servicelevel.NRQL(updateValidFrom),
				Where: servicelevel.NRQL(updateValidWhere),
			}
		}

		if updateGoodFrom != "" {
			eventsInput.GoodEvents = &servicelevel.ServiceLevelEventsQueryUpdateInput{
				From:  servicelevel.NRQL(updateGoodFrom),
				Where: servicelevel.NRQL(updateGoodWhere),
			}
		}

		if updateBadFrom != "" {
			eventsInput.BadEvents = &servicelevel.ServiceLevelEventsQueryUpdateInput{
				From:  servicelevel.NRQL(updateBadFrom),
				Where: servicelevel.NRQL(updateBadWhere),
			}
		}

		input := servicelevel.ServiceLevelIndicatorUpdateInput{
			Name:        updateName,
			Description: updateDescription,
			Events:      &eventsInput,
		}

		result, err := client.NRClient.ServiceLevel.ServiceLevelUpdateWithContext(
			utils.SignalCtx,
			common.EntityGUID(updateGUID),
			updateID,
			input,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdUpdate)

	CmdUpdate.Flags().StringVar(&updateGUID, "guid", "", "The entity GUID of the SLI to update")
	utils.LogIfError(CmdUpdate.MarkFlagRequired("guid"))

	CmdUpdate.Flags().StringVar(&updateID, "id", "", "The unique identifier of the SLI to update")
	utils.LogIfError(CmdUpdate.MarkFlagRequired("id"))

	CmdUpdate.Flags().StringVar(&updateName, "name", "", "The name of the SLI")

	CmdUpdate.Flags().StringVar(&updateDescription, "description", "", "The description of the SLI")

	CmdUpdate.Flags().StringSliceVar(&updateFacets, "facets", []string{}, "The list of NRQL FACET attributes to group results by")

	CmdUpdate.Flags().StringVar(&updateValidFrom, "validEventsFrom", "", "The NRDB event to fetch valid events data from (NRQL FROM clause)")
	CmdUpdate.Flags().StringVar(&updateValidWhere, "validEventsWhere", "", "The NRQL condition to filter valid events")

	CmdUpdate.Flags().StringVar(&updateGoodFrom, "goodEventsFrom", "", "The NRDB event to fetch good events data from (NRQL FROM clause)")
	CmdUpdate.Flags().StringVar(&updateGoodWhere, "goodEventsWhere", "", "The NRQL condition to filter good events")

	CmdUpdate.Flags().StringVar(&updateBadFrom, "badEventsFrom", "", "The NRDB event to fetch bad events data from (NRQL FROM clause)")
	CmdUpdate.Flags().StringVar(&updateBadWhere, "badEventsWhere", "", "The NRQL condition to filter bad events")
}
