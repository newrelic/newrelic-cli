package entities

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
)

var (
	getGUID string
)

var CmdEntitiesGet = &cobra.Command{
	Use:   "get",
	Short: "Get entities",
	Long:  `Get entities via the New Relic API.`,
	Example: `# Get entities
newrelic entities get

# Get with --gUID
newrelic entities get --gUID value`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.Entities.GetEntityWithContext(utils.SignalCtx, common.EntityGUID(getGUID))
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdEntitiesGet)
	CmdEntitiesGet.Flags().StringVar(&getGUID, "gUID", "", "gUID")
	utils.LogIfError(CmdEntitiesGet.MarkFlagRequired("gUID"))
}
