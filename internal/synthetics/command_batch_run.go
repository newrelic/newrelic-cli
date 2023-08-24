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
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	batchFile         string
	guid              []string
	pollingInterval   = time.Second * 30
	progressIndicator = ux.NewSpinner()
	monitorCount      int
	nrdbLatency       = time.Second * 5
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
			config, err = prepareConfig()
			if err != nil {
				log.Fatal(err)
			}
			testsBatchID = runSynthetics(config)

		} else {
			utils.LogIfError(cmd.Help())
		}

		output.Printf("Generated Batch ID: %s", testsBatchID)

		time.Sleep(nrdbLatency)
		handleStatusLoop(accountID, testsBatchID)

	},
}

func init() {
	cmdRun.Flags().StringVarP(&batchFile, "batchFile", "b", "", "Input the YML file to batch and run the monitors")
	cmdRun.Flags().StringSliceVarP(&guid, "guid", "g", nil, "Batch the monitors using their guids and run the automated test")
	Command.AddCommand(cmdRun)

	// MarkFlagsMutuallyExclusive allows one flag at once be invoked
	cmdRun.MarkFlagsMutuallyExclusive("batchFile", "guid")

}

func handleStatusLoop(accountID int, testsBatchID string) {
	// An infinite loop
	ticker := time.NewTicker(pollingInterval)

	defer ticker.Stop()

	for progressIndicator.Start("Fetching the status of tests in the batch...."); true; <-ticker.C {

		// progressIndicator.Start("Fetching the status of tests in the batch....")
		root, err := client.NRClient.Synthetics.GetAutomatedTestResult(accountID, testsBatchID)
		progressIndicator.Stop()

		if err != nil {
			log.Fatal(err)
		}

		log.Println(root.Status, " is the current status")

		exitStatus, ok := TestResultExitCodes[AutomatedTestResultsStatus(root.Status)]

		if !ok {
			exitStatus = handleStatus(*root, AutomatedTestResultsExitStatusUnknown)
		} else {
			exitStatus = handleStatus(*root, exitStatus)
		}

		fmt.Printf("Current Status: %s, Exit Status: %d\n", root.Status, exitStatus)
		os.Stdout.Sync() // Force flush the standard output buffer
		if monitorCount == len(root.Tests) {
			break
		}

	}

}

// getMonitorTestsSummary is called every 15 seconds to print the status of individual monitors
func getMonitorTestsSummary(root synthetics.SyntheticsAutomatedTestResult) (string, [][]string) {
	results := map[string][]synthetics.SyntheticsAutomatedTestJobResult{}

	for _, test := range root.Tests {
		if test.Result == "" {
			test.Result = "PENDING"
		}

		results[string(test.Result)] = append(results[string(test.Result)], test)
	}

	summaryMessage := fmt.Sprintf("%d succeeded; %d failed; %d in progress.",
		len(results["SUCCESS"]), len(results["FAILED"]), len(results["PENDING"]))

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

func handleStatus(root synthetics.SyntheticsAutomatedTestResult, exitStatus AutomatedTestResultsExitStatus) AutomatedTestResultsExitStatus {
	retrievedStatus := string(root.Status)
	switch string(retrievedStatus) {
	case string(AutomatedTestResultsStatusInProgress):
		fmt.Println("\nStatus Received: IN_PROGRESS - re-calling the API in 15 seconds to fetch updated status...")
		summary, tableData := getMonitorTestsSummary(root)
		fmt.Printf("Summary: %s\n", summary)
		printResultTable(tableData)
		return AutomatedTestResultsExitStatusInProgress
	case string(AutomatedTestResultsStatusTimedOut), string(AutomatedTestResultsStatusFailure), string(AutomatedTestResultsStatusPassed):
		progressIndicator.Success("Execution stopped - Status: " + retrievedStatus + "\n")
		fmt.Println("\nStatus Received: " + root.Status + " - Execution halted.")
		summary, tableData := getMonitorTestsSummary(root)
		fmt.Printf("Summary: %s\n", summary)
		printResultTable(tableData)
		return exitStatus
	default:
		progressIndicator.Fail("Unexpected status: " + retrievedStatus)
		fmt.Println("\nStatus Received: " + root.Status + " - Exiting due to unexpected status.")
		return AutomatedTestResultsExitStatusUnknown
	}
}

func printResultTable(tableData [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Status", "Monitor Name", "Monitor GUID", "Is Blocking"})
	table.SetBorder(true) // Set to false to hide the outer border
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	for _, row := range tableData {
		table.Append(row)
	}

	table.Render()
}

// Clean Code

func prepareConfig() (SyntheticsStartAutomatedTestInput, error) {
	if batchFile != "" {
		return readYML(batchFile)
	} else if len(guid) != 0 {
		return createConfigUsingGUIDs(guid), nil
	}
	return SyntheticsStartAutomatedTestInput{}, fmt.Errorf("Invalid arguments")
}

func readYML(batchFile string) (SyntheticsStartAutomatedTestInput, error) {
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

func createConfigUsingGUIDs(guids []string) SyntheticsStartAutomatedTestInput {
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

// runSynthetics batches and call
func runSynthetics(config SyntheticsStartAutomatedTestInput) string {
	log.Println("Batching the following monitors:")
	for _, test := range config.Tests {
		log.Println("-", test.MonitorGUID)
		monitorCount++
	}
	progressIndicator.Start("Batching the monitors")

	result, err := client.NRClient.Synthetics.SyntheticsStartAutomatedTest(config.Config, config.Tests)
	progressIndicator.Stop()
	if err != nil {
		utils.LogIfFatal(err)
	}

	return result.BatchId
}
