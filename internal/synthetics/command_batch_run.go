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
)

var cmdRun = &cobra.Command{
	Use:     "run",
	Example: "newrelic synthetics run --batchFile filename.yml",
	Short:   "Start an automated testing job on a batch of synthetic monitors",
	Long: `Start an automated testing job on a batch of synthetic monitors

The run command helps start an automated testing job by creating a batch, comprising the specified monitors and their
specifications (such as overrides), and subsequently, keeps fetching the status of the batch at periodic intervals of
time, until the status of the batch, which reflects the consolidated status of all monitors in the batch, is either
success, failure or timed out. 

The command may be used with the following flags (the arguments --batchFile and --guid are mutually exclusive).

newrelic synthetics run --batchFile filename.yml
newrelic synthetics run --guid <guid1> --guid <guid2>
`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			config       StartAutomatedTestInput
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
	cmdRun.Flags().StringVarP(&batchFile, "batchFile", "b", "", "Path to the YAML file comprising GUIDs of monitors and associated configuration")
	cmdRun.Flags().StringSliceVarP(&guid, "guid", "g", nil, "List of GUIDs of monitors to include in the batch and run automated tests on")
	Command.AddCommand(cmdRun)

	// MarkFlagsMutuallyExclusive allows one flag at once be invoked
	cmdRun.MarkFlagsMutuallyExclusive("batchFile", "guid")
}

// parseConfiguration helps parse the inputs given to this command, based on the format specified (YAML or command line GUIDs)
func parseConfiguration() (StartAutomatedTestInput, error) {
	if batchFile != "" {
		return createConfigurationUsingYAML(batchFile)
	} else if len(guid) != 0 {
		return createConfigurationUsingGUIDs(guid), nil
	}
	return StartAutomatedTestInput{}, fmt.Errorf("invalid arguments")
}

// createConfigurationUsingYAML unmarshals the specified YAML file into an object that can be used
// to send a create batch request to NerdGraph
func createConfigurationUsingYAML(batchFile string) (StartAutomatedTestInput, error) {
	var config StartAutomatedTestInput

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
func createConfigurationUsingGUIDs(guids []string) StartAutomatedTestInput {
	var tests []synthetics.SyntheticsAutomatedTestMonitorInput
	for _, id := range guids {
		tests = append(tests, synthetics.SyntheticsAutomatedTestMonitorInput{
			MonitorGUID: synthetics.EntityGUID(id),
		})
	}

	return StartAutomatedTestInput{
		Tests: tests,
	}
}

// createAutomatedTestBatch performs an API call to create a batch with the specified configuration and tests
func createAutomatedTestBatch(config StartAutomatedTestInput) string {
	if len(config.Tests) == 0 {
		log.Fatal("No valid monitors found in the input specified. Please check the input provided.")
	}
	progressIndicator.Start("Sending a request to create a batch with the specified monitors...")

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
			if batchResult.Status != synthetics.SyntheticsAutomatedTestStatusTypes.IN_PROGRESS {
				log.Fatal("Unknown Error")
			}
		}

		renderMonitorTestsSummary(*batchResult, exitStatus)

		// Force flush the standard output buffer
		os.Stdout.Sync()

		// exit, if the status is not IN_PROGRESS
		if batchResult.Status != synthetics.SyntheticsAutomatedTestStatusTypes.IN_PROGRESS {
			os.Exit(*exitStatus)
		}
		progressIndicator.Start("Fetching the status of tests in the batch....")
	}
}

// renderMonitorTestsSummary reads through the results of monitors fetched, restructures and renders these results accordingly
func renderMonitorTestsSummary(batchResult synthetics.SyntheticsAutomatedTestResult, exitStatus *int) {
	fmt.Println("Status: ", batchResult.Status, " ")
	summary, tableData := getMonitorTestsSummary(batchResult)
	fmt.Printf("Summary: %s\n", summary)
	if len(tableData) > 0 {
		output.PrintResultTable(tableData)
	}

	if batchResult.Status != synthetics.SyntheticsAutomatedTestStatusTypes.IN_PROGRESS {
		exitStatusMessage := fmt.Sprintf("Exit Status: %d\n", *exitStatus)
		fmt.Println(exitStatusMessage)
	}
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
				string(test.MonitorGUID),
				fmt.Sprintf("%t", test.AutomatedTestMonitorConfig.IsBlocking),
			})
		}
	}

	return summaryMessage, tableData
}
