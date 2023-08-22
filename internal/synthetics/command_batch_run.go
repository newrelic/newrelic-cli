package synthetics

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/newrelic/newrelic-cli/internal/client"
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
		if batchFile != "" || len(guid) != 0 {
			if batchFile != "" {
				config, err = readYML(batchFile)
				if err != nil {
					log.Fatal(err)
				}
			} else if guid != nil {
				// TODO:this is to be tested
				config = createConfigUsingGUIDs(guid)
				log.Println(config)
			}
			testsBatchID = runSynthetics(config)

		} else {
			utils.LogIfError(cmd.Help())
		}

		output.Printf("Generated Batch ID: %s", testsBatchID)

		progressIndicator.Start("Fetching the status of tests in the batch....\n")

		// An infinite loop
		for {
			x := os.Getenv("NEW_RELIC_ACCOUNT_ID")
			y, _ := strconv.Atoi(x)
			root, err := client.NRClient.Synthetics.GetAutomatedTestResult(y, testsBatchID)
			if err != nil {
				log.Fatal(err)
			}

			exitStatus, ok := TestResultExitCodes[AutomatedTestResultsStatus(root.Status)]
			if !ok {
				handleStatus(*root, AutomatedTestResultsExitStatusUnknown)
			} else {
				handleStatus(*root, exitStatus)
			}

		}

	},
}

func init() {
	cmdRun.Flags().StringVarP(&batchFile, "batchFile", "b", "", "Input the YML file to batch and run the monitors")
	cmdRun.Flags().StringSliceVarP(&guid, "guid", "g", nil, "Batch the monitors using their guids and run the automated test")
	Command.AddCommand(cmdRun)

	// MarkFlagsMutuallyExclusive allows one flag at once be invoked
	cmdRun.MarkFlagsMutuallyExclusive("batchFile", "guid")

}

// getMonitorTestsSummary is called every 15 seconds to print the status of individual monitors
func getMonitorTestsSummary(root synthetics.SyntheticsAutomatedTestResult) (string, [][]string) {
    results := map[string]map[string]string{
        "SUCCESS":     {},
        "FAILED":      {},
        "IN_PROGRESS": {},
    }

    for _, test := range root.Tests {
        if test.Result == "" {
            test.Result = "IN_PROGRESS"
        }

        message := fmt.Sprintf("%s (%s)", test.MonitorId, test.MonitorName)
        if test.Result == "FAILED" && test.AutomatedTestMonitorConfig.IsBlocking {
            message = fmt.Sprintf("%s (Blocking)", test.MonitorName)
        }

        results[string(test.Result)][test.MonitorId] = message
    }

    summaryMessage := fmt.Sprintf("%d succeeded; %d failed; %d in progress.",
        len(results["SUCCESS"]), len(results["FAILED"]), len(results["IN_PROGRESS"]))

    tableData := make([][]string, 0)

    for status, messages := range results {
        for _, message := range messages {
            tableData = append(tableData, []string{status, message})
        }
    }

    return summaryMessage, tableData
}

func handleStatus(root synthetics.SyntheticsAutomatedTestResult, exitStatus AutomatedTestResultsExitStatus) {
	retrievedStatus := string(root.Status)
	switch string(retrievedStatus) {
	case string(AutomatedTestResultsStatusInProgress):
		fmt.Println("\nStatus Received: IN_PROGRESS - re-calling the API in 15 seconds to fetch updated status...")
		summary, tableData := getMonitorTestsSummary(root)
		fmt.Printf("Summary: %s\n", summary)
		printResultTable(tableData)
		fmt.Println()
		fmt.Println()
		time.Sleep(pollingInterval)
	case string(AutomatedTestResultsStatusTimedOut), string(AutomatedTestResultsStatusFailure), string(AutomatedTestResultsStatusPassed):
		progressIndicator.Success("Execution stopped - Status: " + retrievedStatus + "\n")
		fmt.Println("\nStatus Received: " + root.Status + " - Execution halted.")
		summary, tableData := getMonitorTestsSummary(root)
		fmt.Printf("Summary: %s\n", summary)
		printResultTable(tableData)
		fmt.Println()
		fmt.Println()
		os.Exit(int(exitStatus))
	default:
		progressIndicator.Fail("Unexpected status: " + retrievedStatus)
		fmt.Println("\nStatus Received: " + root.Status + " - Exiting due to unexpected status.")
		os.Exit(int(AutomatedTestResultsExitStatusUnknown))
	}
}

// Clean Code

func readYML(batchFile string) (SyntheticsStartAutomatedTestInput, error) {
	var config SyntheticsStartAutomatedTestInput
	// Unmarshal YAML file to get monitors and their properties
	// content, err := os.ReadFile(batchFile)
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
	// log.Println("Batching the following monitors:")
	// for _, test := range config.Tests {
	// 	log.Println("-", test.MonitorGUID)
	// }
	progressIndicator.Start("Batching the monitors\n")

	result, err := client.NRClient.Synthetics.SyntheticsStartAutomatedTest(config.Config, config.Tests)
	if err != nil {
		utils.LogIfFatal(err)
	}
	progressIndicator.Stop()
	return result.BatchId
}

func printResultTable(tableData [][]string) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Status", "Monitor"})
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