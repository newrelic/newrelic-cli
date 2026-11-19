package servicelevel

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

var (
	getEntityGUID string
)

var CmdServiceLevelGet = &cobra.Command{
	Use:   "get",
	Short: "Get servicelevel",
	Long:  `Get servicelevel via the New Relic API.`,
	Example: `# Get servicelevel
newrelic servicelevel get

# Get with --entityGuid
newrelic servicelevel get --entityGuid value`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.ServiceLevel.GetIndicatorsWithContext(utils.SignalCtx, common.EntityGUID(getEntityGUID))
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdServiceLevelGet)
	CmdServiceLevelGet.Flags().StringVar(&getEntityGUID, "entityGuid", "", "entityGuid")
	utils.LogIfError(CmdServiceLevelGet.MarkFlagRequired("entityGuid"))
}
