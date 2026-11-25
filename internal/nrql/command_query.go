package nrql

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
)

var (
	historyLimit int
	query        string
)

var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "Execute a NRQL query to New Relic",
	Long: `Execute a NRQL query to New Relic

The query command requires the --query flag which represents a NRQL query string.
This command requires the --accountId <int> flag, which specifies the account to
issue the query against.
`,
	Example: `newrelic nrql query --accountId 12345678 --query 'SELECT count(*) FROM Transaction TIMESERIES'`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()

		_, err := os.Stat(query)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		} else if err == nil {
			// query points to a file/file descriptor
			fileBytes, readErr := os.ReadFile(query)
			if readErr != nil {
				log.Fatal(readErr)
			}
			query = string(fileBytes)
		}

		result, err := client.NRClient.Nrdb.QueryWithContext(utils.SignalCtx, accountID, nrdb.NRQL(query))
		if err != nil {
			log.Fatal(err)
		}

		utils.LogIfFatal(output.Print(result.Results))
	},
}

var cmdHistory = &cobra.Command{
	Use:   "history",
	Short: "Retrieve NRQL query history",
	Long: `Retrieve NRQL query history

The history command will fetch a list of the most recent NRQL queries you executed.
`,
	Example: `newrelic nrql query history`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := client.NRClient.Nrdb.QueryHistory()
		if err != nil {
			log.Fatal(err)
		}

		if result == nil {
			log.Info("no history found. Try using the 'newrelc nrql query' command")
			return
		}

		count := len(*result)

		if count < historyLimit {
			historyLimit = count
		}

		output.Text((*result)[0:historyLimit])
	},
}

func init() {
	Command.AddCommand(cmdQuery)

	cmdQuery.Flags().StringVarP(&query, "query", "q", "", "the NRQL query you want to execute")
	utils.LogIfError(cmdQuery.MarkFlagRequired("query"))

	Command.AddCommand(cmdHistory)
	cmdHistory.Flags().IntVarP(&historyLimit, "limit", "l", 10, "history items to return (default: 10, max: 100)")
}
