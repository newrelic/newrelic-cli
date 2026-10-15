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

var CmdWorkloadsGet = &cobra.Command{
	Use:   "get",
	Short: "Get workloads",
	Long:  `Get workloads via the New Relic API.`,
	Example: `# Get workloads
newrelic workload get

# Get with --accountID
newrelic workload get --accountID 123`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.Workloads.GetCollectionWithContext(utils.SignalCtx, getAccountID, common.EntityGUID(getGUID))
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdWorkloadsGet)
	CmdWorkloadsGet.Flags().IntVar(&getAccountID, "accountID", 0, "accountID")
	CmdWorkloadsGet.Flags().StringVar(&getGUID, "gUID", "", "gUID")
	utils.LogIfError(CmdWorkloadsGet.MarkFlagRequired("accountID"))
	utils.LogIfError(CmdWorkloadsGet.MarkFlagRequired("gUID"))
}
