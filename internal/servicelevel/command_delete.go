package servicelevel

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

var (
	deleteGUID string
)

var CmdServiceLevelDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete servicelevel",
	Long:  `Delete servicelevel via the New Relic API.`,
	Example: `# Delete servicelevel
newrelic servicelevel delete

# Delete with --guid
newrelic servicelevel delete --guid value`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.ServiceLevel.ServiceLevelDeleteWithContext(utils.SignalCtx, common.EntityGUID(deleteGUID))
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdServiceLevelDelete)
	CmdServiceLevelDelete.Flags().StringVar(&deleteGUID, "guid", "", "guid")
	utils.LogIfError(CmdServiceLevelDelete.MarkFlagRequired("guid"))
}
