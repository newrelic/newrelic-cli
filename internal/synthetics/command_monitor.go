package synthetics

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

var (
	statusFilter string
	monitorName  string
	monitorID    string
)

// Command represents the synthetics command
var cmdMon = &cobra.Command{
	Use:     "monitor",
	Short:   "Interact with New Relic Synthetics monitors",
	Example: "newrelic synthetics monitor --help",
	Long:    "Interact with New Relic Synthetics monitors",
}

var cmdMonGet = &cobra.Command{
	Use:   "get",
	Short: "Get a New Relic Synthetics monitor",
	Long: `Get a New Relic Synthetics monitor

The get command performs a query for an Synthetics monitor by ID.
`,
	Example: `newrelic synthetics monitor get --monitorId "<monitorID>"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var results *synthetics.Monitor
		var err error

		if monitorID != "" {
			results, err = client.NRClient.Synthetics.GetMonitor(monitorID)
			utils.LogIfFatal(err)
		} else {
			utils.LogIfError(cmd.Help())
			log.Fatal(" --monitorId <monitorID> is required")
		}

		utils.LogIfFatal(output.Print(results))
	},
}

var cmdMonList = &cobra.Command{
	Use:   "list",
	Short: "List New Relic Synthetics monitors",
	Long: `List New Relic Synthetics monitors

The list command performs a query for all Synthetics monitors, optionally filtered on the status field.
`,
	Example: `newrelic synthetics monitor list --statusFilter "DISABLED, MUTED"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var results []*synthetics.Monitor
		var err error
		var jqFilter string

		results, err = client.NRClient.Synthetics.ListMonitors()
		utils.LogIfFatal(err)

		if statusFilter != "" {
			statusFilter = strings.ToUpper(statusFilter)
			var replacer = strings.NewReplacer(
				" ", "",
				"\r\n", "",
				"\n\r", "",
				"\r", "",
				"\n", "",
				"\t", "",
			)
			statusFilter = replacer.Replace(statusFilter)
			statusFilter = strings.ReplaceAll(statusFilter, ",", `", "`)
			jqFilter = `.[] | select(.status as $s | ["` + statusFilter + `"] | index($s))`
			query, err := gojq.Parse(jqFilter)
			if err != nil {
				utils.LogIfFatal(err)
			}

			bytes, err := json.Marshal(results)
			if err != nil {
				utils.LogIfFatal(err)
			}

			var obj interface{}
			err = json.Unmarshal(bytes, &obj)
			if err != nil {
				utils.LogIfFatal(err)
			}

			iter := query.Run(obj)
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}

				if err, ok := v.(error); ok {
					utils.LogIfFatal(err)
				}

				utils.LogIfFatal(output.Print(v))
			}
		} else {
			utils.LogIfFatal(output.Print(results))
		}
	},
}

var cmdMonSearch = &cobra.Command{
	Use:   "search",
	Short: "Search for a New Relic Synthetics Monitor",
	Long: `Search for a New Relic Synthetics Monitor

The search command performs a query for a Synthetics Monitor by name.
`,
	Example: "newrelic synthetics monitor search --name <monitorName>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.GetActiveProfileAccountID()

		if monitorName == "" && accountID == 0 {
			utils.LogIfError(cmd.Help())
			log.Fatal("one of --accountId or --name are required")
		}

		var entityResults []entities.EntityOutlineInterface
		var err error

		params := entities.EntitySearchQueryBuilder{
			Domain: entities.EntitySearchQueryBuilderDomain("SYNTH"),
			Type:   entities.EntitySearchQueryBuilderType("MONITOR"),
		}

		if monitorName != "" {
			params.Name = monitorName
		}

		if accountID != 0 {
			params.Tags = []entities.EntitySearchQueryBuilderTag{{Key: "accountId", Value: strconv.Itoa(accountID)}}
		}

		results, err := client.NRClient.Entities.GetEntitySearch(
			entities.EntitySearchOptions{},
			"",
			params,
			[]entities.EntitySearchSortCriteria{},
		)
		if err != nil {
			utils.LogIfFatal(err)
		}

		entityResults = results.Results.Entities
		utils.LogIfFatal(output.Print(entityResults))

	},
}

func init() {
	Command.AddCommand(cmdMon)

	cmdMonGet.Flags().StringVarP(&monitorID, "monitorId", "", "", "A New Relic Synthetics monitor ID")
	cmdMon.AddCommand(cmdMonGet)

	cmdMonList.Flags().StringVarP(&statusFilter, "statusFilter", "s", "", "filter the results on the status field. Possible values ENABLED, DISABLED, MUTED. Comma separated.")
	cmdMon.AddCommand(cmdMonList)

	cmdMonSearch.Flags().StringVarP(&monitorName, "name", "n", "", "search for results matching the given Synthetics monitor name")
	cmdMon.AddCommand(cmdMonSearch)
}
