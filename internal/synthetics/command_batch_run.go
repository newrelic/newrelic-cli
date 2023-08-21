package synthetics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	batchFile string
	guid      []string

	// Temporary variable to shift between scenarios of
	// mock JSON responses of the automatedTestResults query
	scenario = 2

	apiTimeout = time.Minute * 5

	// To be changed to 30 seconds in the real implementation, as suggested by the Synthetics team
	pollingInterval = time.Second * 15

	// Spinner
	//progressIndicator = ux.NewSpinnerProgressIndicator()
	progressIndicator = ux.NewSpinner()
)

var cmdRun = &cobra.Command{
	Use: "run",
	//TODO: Find the Precise description.
	Short:   "Interact with New Relic Synthetics batch monitors",
	Example: "newrelic synthetics run --help",
	Long:    "Interact with New Relic Synthetics batch monitors",
	Run: func(cmd *cobra.Command, args []string) {
		var mockbatchID string

		// Config holds values unmarshalled from the YAML file
		var config StartAutomatedTestInput

		if batchFile != "" || len(guid) != 0 {
			if batchFile != "" {
				// Unmarshal YAML file to get monitors and their properties
				// content, err := os.ReadFile(batchFile)
				content, err := os.ReadFile(batchFile)
				if err != nil {
					log.Fatal(err)
				}
				err = yaml.Unmarshal(content, &config)
				if err != nil {
					log.Fatal(err)
				}

				for _, test := range config.Tests {
					requestBody := fmt.Sprintf(`{"guid": "%s", "isBlocking": %v}`, test.MonitorGuid, test.Config.IsBlocking)
					fmt.Println(requestBody)

					// appending GUIDs from the YAML to []guid as the arguments --guid and --batchFile
					// are mutually exclusive - so the same variable can hold guids from the YAML.
					guid = append(guid, test.MonitorGuid)
				}

			} else if guid != nil {
				// do nothing, as guid is already populated with --guid arguments given in the command
			}

			// ----------------------------------------------------------------------------------
			// NOTE: mockBatchID is returned by the function that is (in reality) expected to
			// typecast inputs as needed by the API, send a request, and receive a batchID
			// as the response. Since this is being mocked, we only send the list of guids
			// to the function for now, but this is expected to belong to newrelic-client-go
			// and take all required parameters from the input YAML or GUIDs.
			// ----------------------------------------------------------------------------------

			mockbatchID = runSynthetics(guid)

		} else {
			utils.LogIfError(cmd.Help())
			//log.Fatal(" --batchFile <ymlFile> is required")
		}

		utils.LogIfFatal(output.Print(mockbatchID))

		// ----------------------------------------------------------------------------------
		// In order to mock implementation, the batchID has been hardcoded.
		// This is expected to be received in the response of syntheticsStartAutomatedTest.
		// ----------------------------------------------------------------------------------

		progressIndicator.Start("Fetching the status of the batch\n")
		start := time.Now()

		// This variable has been used to track iterations of mock API calls for easier file opening of mock JSON files.
		// This may be discarded along with its usages after the API calling function is used.
		i := 0

		// An infinite loop
		for {

			// A timeout in the CLI may not be needed, based on recent suggestions received, as the API
			// returns a TIMED_OUT status if one or more job(s) in the batch consume > 10 minutes.
			if time.Since(start) > apiTimeout {
				fmt.Println("---------------------------")
				progressIndicator.Fail("Halting execution : reached timeout.")
				fmt.Println("---------------------------")
				break
			}

			i++

			root := getAutomatedTestResultsMockResponse(mockbatchID, i)
			// ----------------------------------------------------------------------------------
			// ORIGINAL FUNCTION : Call the method from go-client that would send a request to the
			// automatedTestResults query with the batchID and fetch monitor details - sample below.

			// result, err := functionInClientGo(automatedTestResultQueryInput{batchID: batchID})
			// if err != nil {
			//	 return fmt.Errorf("Some error")
			// }
			// ----------------------------------------------------------------------------------

			exitStatus, ok := TestResultExitCodes[AutomatedTestResultsStatus(root.Status)]
			if !ok {
				handleStatus(root, AutomatedTestResultsExitStatusUnknown)
			} else {
				handleStatus(root, exitStatus)
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

// getAutomatedTestResultsMockResponse is called to retrieve the mock JSON response of the automatedTestResults query
func getAutomatedTestResultsMockResponse(batchID string, index int) (r Root) {
	directory := fmt.Sprintf("internal/synthetics/mock_json/Scenario %d", scenario)
	filePath := fmt.Sprintf("%s/response_%d.json", directory, index)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(jsonFile *os.File) {
		err = jsonFile.Close()
		if err != nil {
			log.Fatal("Unable to close the file")
		}
	}(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var root Root
	if err := json.Unmarshal(byteValue, &root); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	return root
}

// printMonitorStatus is called every 15 seconds to print the status of individual monitors
func printMonitorStatus(root Root) {
	countSuccess := 0
	countFailure := 0
	countProgress := 0
	tests := root.Tests

	for _, test := range tests {
		if test.Result == "SUCCESS" {
			countSuccess++
		} else if test.Result == "FAILED" {
			countFailure++
		} else if test.Result == "IN_PROGRESS" || test.Result == "" {
			countProgress++
		}
	}

	fmt.Println("Successful Tests: ", countSuccess)
	fmt.Println("Failed Tests: ", countFailure)
}

// runSynthetics batches and call
func runSynthetics(guids []string) string {
	fmt.Println("Running New Relic Synthetics for the following monitor GUID(s):")
	for _, guid := range guids {
		fmt.Println("-", guid)
	}

	// ----------------------------------------------------------------------------------
	// ORIGINAL FUNCTION : Call the method from go-client that would send a request to the
	// syntheticsStartAutomatedTest mutation (with the request body fit into datatypes
	// generated by Tutone) and would receive a response comprising a batchID.
	// ----------------------------------------------------------------------------------

	// returning a mock batchID.
	// Will be replaced with the response from the API
	return "36488dff-9a8a-4358-8ef2-e73da3c118e0"
}

// handleStatus processes the execution result of a test, taking into account
// the status contained in the root object and the associated exit status.
// Depending on the root.Status value, it prints an appropriate message, waits
// for the next API call, or exits the program with the given exit status.
//
// Parameters:
//   - root: Root struct that contains the status information.
//   - exitStatus: The AutomatedTestResultsExitStatus corresponding to the given root.Status.
//
// In the case of AutomatedTestResultsStatusInProgress, the function prints an
// information message, calls the printMonitorStatus function, and waits for
// the specified pollingInterval before the next API call.
//
// In the cases of AutomatedTestResultsStatusTimedOut, AutomatedTestResultsStatusFailure,
// and AutomatedTestResultsStatusPassed, the function prints the execution result,
// calls the printMonitorStatus function, and exits the program with the
// corresponding exit status code.
func handleStatus(root Root, exitStatus AutomatedTestResultsExitStatus) {
	switch root.Status {
	case string(AutomatedTestResultsStatusInProgress):
		fmt.Println("Status still IN_PROGRESS, calling API again in 15 seconds")
		printMonitorStatus(root)
		fmt.Println()
		time.Sleep(pollingInterval)
	case string(AutomatedTestResultsStatusTimedOut), string(AutomatedTestResultsStatusFailure), string(AutomatedTestResultsStatusPassed):
		progressIndicator.Success("Execution stopped - Status: " + root.Status + "\n")
		printMonitorStatus(root)
		os.Exit(int(exitStatus))
	default:
		progressIndicator.Fail("Unexpected status: " + root.Status + "\n")
		os.Exit(int(AutomatedTestResultsExitStatusUnknown))
	}
}
