package synthetics

import (
	"fmt"
	"os"
	"time"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	batchFile         string
	guid              []string
	pollingInterval   = time.Second * 30
	progressIndicator = ux.NewSpinner()
	nrdbLatency       = time.Second * 5
	countMonitors     int
)

var cmdRun = &cobra.Command{
	Use: "run",
	//TODO: Find the Precise description.
	Short:   "Interact with New Relic Synthetics batch monitors",
	Example: "newrelic synthetics run --help",
	Long:    "Interact with New Relic Synthetics batch monitors",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			config       SyntheticsStartAutomatedTestInput
			err          error
			testsBatchID string
		)
		accountID := configAPI.GetActiveProfileAccountID()
		if batchFile != "" || len(guid) != 0 {
			config, err = parseConfiguration()
			if err != nil {
				log.Fatal(err)
			}

			testsBatchID = createAutomatedTestBatch(config)
			output.Printf("Generated Batch ID: %s", testsBatchID)

			// can be ignored if there is no initial tick by the ticker
			time.Sleep(nrdbLatency)
			getAutomatedTestResults(accountID, testsBatchID)

		} else {
			utils.LogIfError(cmd.Help())
		}

	},
}

// Definition of the command
func init() {
	cmdRun.Flags().StringVarP(&batchFile, "batchFile", "b", "", "Input the YML file to batch and run the monitors")
	cmdRun.Flags().StringSliceVarP(&guid, "guid", "g", nil, "Batch the monitors using their guids and run the automated test")
	Command.AddCommand(cmdRun)

	// MarkFlagsMutuallyExclusive allows one flag at once be invoked
	cmdRun.MarkFlagsMutuallyExclusive("batchFile", "guid")
}

// parseConfiguration helps parse the inputs given to this command, based on the format specified (YAML or command line GUIDs)
func parseConfiguration() (SyntheticsStartAutomatedTestInput, error) {
	if batchFile != "" {
		return createConfigurationUsingYAML(batchFile)
	} else if len(guid) != 0 {
		return createConfigurationUsingGUIDs(guid), nil
	}
	return SyntheticsStartAutomatedTestInput{}, fmt.Errorf("Invalid arguments")
}

// createConfigurationUsingYAML unmarshals the specified YAML file into an object that can be used
// to send a create batch request to NerdGraph
func createConfigurationUsingYAML(batchFile string) (SyntheticsStartAutomatedTestInput, error) {
	var config SyntheticsStartAutomatedTestInput

	content, err := os.ReadFile(batchFile)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	utils.LogIfFatal(err)
	return config, nil
}

// createConfigurationUsingGUIDs obtains GUIDs specified in command line arguments and restructures them into an object
// that can be used to send a create batch request to NerdGraph
func createConfigurationUsingGUIDs(guids []string) SyntheticsStartAutomatedTestInput {
	var tests []synthetics.SyntheticsAutomatedTestMonitorInput
	for _, id := range guids {
		tests = append(tests, synthetics.SyntheticsAutomatedTestMonitorInput{
			MonitorGUID: synthetics.EntityGUID(id),
		})
	}

	return SyntheticsStartAutomatedTestInput{
		Tests: tests,
	}
}

// createAutomatedTestBatch performs an API call to create a batch with the specified configuration and tests
func createAutomatedTestBatch(config SyntheticsStartAutomatedTestInput) string {
	log.Println("Creating a batch comprising the following monitors:")
	for _, test := range config.Tests {
		log.Println("-", test.MonitorGUID)
		countMonitors++
	}

	if countMonitors == 0 {
		log.Fatal("No valid monitors found in the input specified. Please check the input provided.")
	}
	progressIndicator.Start("Sending a request to create the batch:")

	result, err := client.NRClient.Synthetics.SyntheticsStartAutomatedTest(config.Config, config.Tests)
	progressIndicator.Stop()
	if err != nil {
		utils.LogIfFatal(err)
	}

	return result.BatchId
}

// getAutomatedTestResults performs an API call at regular intervals of time (when the pollingInterval has elapsed)
// to fetch the consolidated status of the batch, and the results of monitors the batch comprises
func getAutomatedTestResults(accountID int, testsBatchID string) {
	// An infinite loop
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for progressIndicator.Start("Fetching the status of tests in the batch...."); true; <-ticker.C {
		batchResult, err := client.NRClient.Synthetics.GetAutomatedTestResult(accountID, testsBatchID)
		progressIndicator.Stop()

		if err != nil {
			log.Fatal(err)
		}

		exitStatus, ok := globalResultExitCodes[(batchResult.Status)]
		if !ok {
			log.Fatal("Unknown Error")
		} else {
			renderMonitorTestsSummary(*batchResult, exitStatus)
		}

		fmt.Printf("Current Status: %s, Exit Status: %d\n", batchResult.Status, *exitStatus)

		// Force flush the standard output buffer
		os.Stdout.Sync()

		// exit, if the status is not IN_PROGRESS
		if batchResult.Status != synthetics.SyntheticsAutomatedTestStatusTypes.IN_PROGRESS {
			break
		}
		progressIndicator.Start("Fetching the status of tests in the batch....")
	}
}

// renderMonitorTestsSummary reads through the results of monitors fetched, restructures and renders these results accordingly
func renderMonitorTestsSummary(batchResult synthetics.SyntheticsAutomatedTestResult, exitStatus *int) {
	fmt.Println("Status Received: ", batchResult.Status, " ")
	summary, tableData := getMonitorTestsSummary(batchResult)
	fmt.Printf("Summary: %s\n", summary)
	printResultTable(tableData)
}

// getMonitorTestsSummary reads through the results of monitors fetched and populates them to a table with details
// of each monitor, to print these results to the terminal
func getMonitorTestsSummary(batchResult synthetics.SyntheticsAutomatedTestResult) (string, [][]string) {
	results := map[string][]synthetics.SyntheticsAutomatedTestJobResult{}

	for _, test := range batchResult.Tests {
		if test.Result == "" {
			test.Result = synthetics.SyntheticsJobStatusTypes.PENDING
		}

		results[string(test.Result)] = append(results[string(test.Result)], test)
	}

	summaryMessage := fmt.Sprintf("%d succeeded; %d failed; %d in progress.",
		len(results[string(synthetics.SyntheticsJobStatusTypes.SUCCESS)]),
		len(results[string(synthetics.SyntheticsJobStatusTypes.FAILED)]),
		len(results[string(synthetics.SyntheticsJobStatusTypes.PENDING)]))

	tableData := make([][]string, 0)

	for status, tests := range results {
		for _, test := range tests {
			tableData = append(tableData, []string{
				status,
				test.MonitorName,
				string(test.MonitorId),
				fmt.Sprintf("%t", test.AutomatedTestMonitorConfig.IsBlocking),
			})
		}
	}

	return summaryMessage, tableData
}
