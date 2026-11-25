package workload

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

var (
	getAccountID int
	getGUID      string
)

var CmdWorkloadsGetCollection = &cobra.Command{
	Use:   "get-collection",
	Short: "Get workloads",
	Long:  `Get workloads via the New Relic API.`,
	Example: `# Get workloads
newrelic workload get

# Get with --accountId
newrelic workload get --accountId 123`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.Workloads.GetCollectionWithContext(utils.SignalCtx, getAccountID, common.EntityGUID(getGUID))
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdWorkloadsGetCollection)
	CmdWorkloadsGetCollection.Flags().IntVar(&getAccountID, "accountId", 0, "accountId")
	CmdWorkloadsGetCollection.Flags().StringVar(&getGUID, "guid", "", "guid")
	utils.LogIfError(CmdWorkloadsGetCollection.MarkFlagRequired("accountId"))
	utils.LogIfError(CmdWorkloadsGetCollection.MarkFlagRequired("guid"))
}
