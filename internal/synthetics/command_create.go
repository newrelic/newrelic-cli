package synthetics

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	createAccountID               int
	createDescription             string
	createName                    string
	createShared                  bool
	createVerifiedScriptExecution bool
)

var CmdSyntheticsCreate = &cobra.Command{
	Use:   "create",
	Short: "Create synthetics",
	Long:  `Create synthetics via the New Relic API.`,
	Example: `# Create synthetics
newrelic synthetics create

# Create with --accountId
newrelic synthetics create --accountId 123`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.Synthetics.SyntheticsCreatePrivateLocationWithContext(utils.SignalCtx, createAccountID, createDescription, createName, createShared, createVerifiedScriptExecution)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdSyntheticsCreate)
	CmdSyntheticsCreate.Flags().IntVar(&createAccountID, "accountId", 0, "accountId")
	CmdSyntheticsCreate.Flags().StringVar(&createDescription, "description", "", "description")
	CmdSyntheticsCreate.Flags().StringVar(&createName, "name", "", "name")
	CmdSyntheticsCreate.Flags().BoolVar(&createShared, "shared", false, "shared")
	CmdSyntheticsCreate.Flags().BoolVar(&createVerifiedScriptExecution, "verifiedScriptExecution", false, "verifiedScriptExecution")
	utils.LogIfError(CmdSyntheticsCreate.MarkFlagRequired("accountId"))
	utils.LogIfError(CmdSyntheticsCreate.MarkFlagRequired("description"))
	utils.LogIfError(CmdSyntheticsCreate.MarkFlagRequired("name"))
}
