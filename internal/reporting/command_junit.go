package reporting

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/joshdk/go-junit"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

const junitEventType = "TestRun"

var (
	path         string
	attributes   string
	dryRun       bool
	outputEvents bool
)

var cmdJUnit = &cobra.Command{
	Use:   "junit",
	Short: "Send JUnit test run results to New Relic",
	Long: `Send JUnit test run results to New Relic

`,
	Example: `newrelic reporting junit --accountId 12345678 --path unit.xml --attributes '{"sha": 12345}'`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()

		if configAPI.GetActiveProfileString(config.LicenseKey) == "" {
			log.Fatal("a License key is required, set one in your default profile or use the NEW_RELIC_LICENSE_KEY environment variable")
		}

		id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}

		xml, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		suites, err := junit.Ingest(xml)
		if err != nil {
			log.Fatalf("failed to ingest JUnit xml %v", err)
		}

		var a map[string]interface{}

		err = json.Unmarshal([]byte(attributes), &a)
		if err != nil {
			log.Fatal(err)
		}

		events := []map[string]interface{}{}

		for _, suite := range suites {
			for _, test := range suite.Tests {
				events = append(events, createTestRunEvent(id, suite, test, a))
			}
		}

		if outputEvents {
			utils.LogIfFatal(output.Print(events))
		}

		if dryRun {
			return
		}

		if err := client.NRClient.Events.CreateEventWithContext(utils.SignalCtx, accountID, events); err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

func createTestRunEvent(testRunID uuid.UUID, suite junit.Suite, test junit.Test, attributes map[string]interface{}) map[string]interface{} {
	e := map[string]interface{}{}
	e["eventType"] = junitEventType
	e["id"] = testRunID.String()
	e["test"] = test.Name
	e["classname"] = test.Classname
	e["suite"] = suite.Name
	e["package"] = suite.Package
	e["status"] = test.Status
	e["durationMs"] = test.Duration.Milliseconds()

	if test.Error != nil {
		e["errorMessage"] = test.Error.Error()
	}

	for key, value := range suite.Properties {
		e[key] = value
	}

	for key, value := range test.Properties {
		e[key] = value
	}

	for key, value := range attributes {
		e[key] = value

	}

	return e
}

func init() {
	Command.AddCommand(cmdJUnit)
	cmdJUnit.Flags().StringVarP(&path, "path", "p", "", "the path to a JUnit-formatted test results file")
	cmdJUnit.Flags().StringVarP(&attributes, "attributes", "", "{}", "any custom attributes to include")
	cmdJUnit.Flags().BoolVarP(&outputEvents, "output", "o", false, "output generated custom events to stdout")
	cmdJUnit.Flags().BoolVar(&dryRun, "dryRun", false, "suppress posting custom events to NRDB")
	utils.LogIfError(cmdJUnit.MarkFlagRequired("path"))
}
