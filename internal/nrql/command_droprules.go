package nrql

import (
	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dropRulesLimit int
)

var cmdDropRules = &cobra.Command{
	Use:   "droprules",
	Short: "Retrieve NRQL Drop Rules",
	Long: `Retrieve NRQL Drop Rules

The 'droprules' command helps fetch NRQL droprules associated with your account.
`,
	Example: `newrelic nrql droprules --limit <limit>`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()
		result, err := client.NRClient.Nrqldroprules.GetList(accountID)
		if err != nil {
			log.Fatal(err)
		}

		if result == nil {
			log.Info("No drop rules found, associated with your account.")
			return
		}
		rulesResult := *result
		dropRules := rulesResult.Rules
		lengthOfDropRules := len(dropRules)

		if lengthOfDropRules == 0 {
			log.Info("No drop rules found, associated with your account.")
			return
		}

		if dropRulesLimit == 0 {
			log.Info("A minimum of '1' needs to be specified with the --limit flag. If not specified, the limit defaults to 10.")
			return
		}

		if lengthOfDropRules < dropRulesLimit {
			dropRulesLimit = lengthOfDropRules
		}
		dropRulesSliced := (dropRules)[0:dropRulesLimit]

		outputErr := output.Print(dropRulesSliced)
		if outputErr != nil {
			utils.LogIfFatal(err)
		}

	},
}

func init() {
	Command.AddCommand(cmdDropRules)
	cmdDropRules.Flags().IntVarP(&dropRulesLimit, "limit", "l", 10, "Number of NRQL Drop Rules to return")
}
