package apm

import (
	"fmt"
	"log"

	"github.com/hokaccha/go-prettyjson"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/spf13/cobra"
)

var (
	nrClient         *newrelic.NewRelic
	apmApplicationID int
)

// SetClient is the API for passing along the New Relic client to this command
func SetClient(nr *newrelic.NewRelic) error {
	if nr == nil {
		return fmt.Errorf("client can not be nil")
	}

	nrClient = nr

	return nil
}

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "apm",
	Short: "apm",
}

var apmDescribeDeployments = &cobra.Command{
	Use:   "describe-deployments",
	Short: "describe-deployments",
	Long: `Search for New Relic APM deployments

The describe-deployments command performs a search for New Relic APM
deployments.
`,
	Example: "newrelic apm describe-deployments --applicationID <appID>",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		fmt.Println(nrClient)

		deployments, err := nrClient.APM.ListDeployments(apmApplicationID)
		if err != nil {
			log.Fatal(err)
		}

		json, err := prettyjson.Marshal(deployments)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	},
}

func init() {
	Command.AddCommand(apmDescribeDeployments)
	apmDescribeDeployments.Flags().IntVarP(&apmApplicationID, "applicationID", "a", 0, "search for results matching the given name")
	apmDescribeDeployments.MarkFlagRequired("applicationID")
}
