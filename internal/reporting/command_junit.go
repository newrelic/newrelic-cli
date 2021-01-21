package reporting

import (
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/joshdk/go-junit"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/output"
)

const junitEventType = "TestRun"

var (
	accountID    int
	path         string
	dryRun       bool
	outputEvents bool
)

var cmdJUnit = &cobra.Command{
	Use:   "junit",
	Short: "Send JUnit test run results to New Relic",
	Long: `Send JUnit test run results to New Relic

`,
	Example: `newrelic reporting junit --accountId 12345678 --path unit.xml`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if accountID, err = config.RequireAccountID(); err != nil {
			log.Fatal(err)
		}

		if _, err = config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		insightsInsertKey := config.GetActiveProfileValueString(config.InsightsInsertKey)
		if insightsInsertKey == "" {
			log.Fatal("an Insights insert key is required, set one in your default profile or use the NEW_RELIC_INSIGHTS_INSERT_KEY environment variable")
		}

		id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}

		xml, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		suites, err := junit.Ingest(xml)
		if err != nil {
			log.Fatalf("failed to ingest JUnit xml %v", err)
		}

		events := []map[string]interface{}{}

		for _, suite := range suites {
			for _, test := range suite.Tests {
				events = append(events, createTestRunEvent(id, suite, test))
			}
		}

		if outputEvents {
			if err := output.Print(events); err != nil {
				log.Fatal(err)
			}
		}

		if dryRun {
			return
		}

		if err := client.Client.Events.CreateEvent(accountID, events); err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

func createTestRunEvent(testRunID uuid.UUID, suite junit.Suite, test junit.Test) map[string]interface{} {
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

	return e
}

func init() {
	Command.AddCommand(cmdJUnit)
	cmdJUnit.Flags().StringVarP(&path, "path", "p", "", "the path to a JUnit-formatted test results file")
	cmdJUnit.Flags().BoolVarP(&outputEvents, "output", "o", false, "output generated custom events to stdout")
	cmdJUnit.Flags().BoolVar(&dryRun, "dryRun", false, "suppress posting custom events to NRDB")
	if err := cmdJUnit.MarkFlagRequired("path"); err != nil {
		log.Error(err)
	}
}
