package nerdgraph

import (
	"bytes"
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/output"
	ng "github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
)

var (
	variables string
)

var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "Execute a raw GraphQL query request to the NerdGraph API",
	Long: `Execute a raw GraphQL query request to the NerdGraph API

The query command accepts a single argument in the form of a GraphQL query as a string.
This command accepts an optional flag, --variables, which should be a JSON string where the
keys are the variables to be referenced in the GraphQL query.
`,
	Example: `newrelic nerdgraph query 'query($guid: EntityGuid!) { actor { entity(guid: $guid) { guid name domain entityType } } }' --variables '{"guid": "<GUID>"}'`,
	Args: func(cmd *cobra.Command, args []string) error {
		argsCount := len(args)

		if argsCount < 1 {
			return errors.New("missing graph query argument")
		}

		if argsCount > 1 {
			return errors.New("command expects only 1 argument")
		}

		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var variablesParsed map[string]interface{}

		if err := json.Unmarshal([]byte(variables), &variablesParsed); err != nil {
			log.Fatal(err)
		}

		query := args[0]

		result, err := client.Client.NerdGraph.Query(query, variablesParsed)
		if err != nil {
			log.Fatal(err)
		}

		reqBodyBytes := new(bytes.Buffer)

		encoder := json.NewEncoder(reqBodyBytes)
		if err = encoder.Encode(ng.QueryResponse{
			Actor: result.(ng.QueryResponse).Actor,
		}); err != nil {
			log.Fatal(err)
		}

		if err := output.Print(reqBodyBytes); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdQuery)
	cmdQuery.Flags().StringVar(&variables, "variables", "{}", "the variables to pass to the GraphQL query, represented as a JSON string")
}
