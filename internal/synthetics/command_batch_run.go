package synthetics

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"
)

// MockClient is a mock implementation of the HTTPClient interface using testify/mock.
type MockClient struct {
	mock.Mock
}

// Do is the mocked version of the http.Client Do method.
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

var (
	batchFile string
)

var scenario = 2

// Command represents the synthetics command
var cmdRun = &cobra.Command{
	Use: "run",
	//TODO: Find the Precise description.
	Short:   "Interact with New Relic Synthetics batch monitors",
	Example: "newrelic synthetics run --help",
	Long:    "Interact with New Relic Synthetics monitors",
	//TODO: Start working on the mocking the json
	Run: func(cmd *cobra.Command, args []string) {
		var results *synthetics.Monitor
		// Config holds the value of the YML
		var config Configuration

		// Create the mock HTTP client
		mockClient := &MockClient{}

		if batchFile != "" {
			// YAML file input Unmarshall from a file
			content, err := os.ReadFile(batchFile)
			if err != nil {
				log.Fatal(err)
			}
			err = yaml.Unmarshal(content, &config)
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println(batchFile)
			//fmt.Println(config.Monitors)

			for _, monitor := range config.Monitors {
				apiURL := "https://example.com/api"
				requestBody := fmt.Sprintf(`{"guid": "%s", "isBlocking": %v}`, monitor.GUID, monitor.Config.IsBlocking)
				fmt.Println(requestBody)

				// Create the request
				req, err := http.NewRequest("POST", apiURL, strings.NewReader(requestBody))
				if err != nil {
					fmt.Println("Error creating request:", err)
					return
				}

				// Expect a response for the request
				expectedResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader(`{"status": "success", "message": "Mock response from API"}`)),
				}

				mockClient.On("Do", req).Return(expectedResponse, nil)

				// Perform the API call using the mockClient
				resp, err := mockClient.Do(req)
				if err != nil {
					fmt.Println("API call error:", err)
					return
				}

				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						log.Fatal(err)
					}
				}(resp.Body)

				// fmt.Println("Printing response")
				// fmt.Println(resp)
			}

			// TODO: replace with a mock function
			//results, err = client.NRClient.Synthetics.GetMonitor(batchFile)
			//utils.LogIfFatal(err)
		} else {
			utils.LogIfError(cmd.Help())
			log.Fatal(" --batchFile <ymlFile> is required")
		}

		utils.LogIfFatal(output.Print(results))

		batchID := "36488dff-9a8a-4358-8ef2-e73da3c118e0"
		apiTimeout := time.Minute * 5

		start := time.Now()
		i := 0
		for {
			if time.Since(start) > apiTimeout {
				fmt.Println("Halting execution : reached timeout.")
				break
			}

			i++

			//result, err := functionInClientGo(automatedTestResultQueryInput{batchID: batchID})
			root := fakeAutomatedTestResultQuery(batchID, i)
			//if err != nil {
			//	return fmt.Errorf("Some error")
			//}

			if root.Status == "TIMED_OUT" || root.Status == "PASSED" || root.Status == "FAILED" {
				fmt.Printf("Execution stopped - Status: %s\n", root.Status)
				fakePrintMonitorStatus(root)
				break
			} else if root.Status == "IN_PROGRESS" {
				fmt.Println("Status still IN_PROGRESS, calling API again in 15 seconds")
				fakePrintMonitorStatus(root)
				fmt.Println()
				time.Sleep(time.Second * 15)
			} else {
				fmt.Printf("Unexpected status: %s\n", root.Status)
				break
			}
		}

	},
}

func fakeAutomatedTestResultQuery(batchID string, index int) (r Root) {
	directory := fmt.Sprintf("internal/synthetics/mock_json/Scenario %d", scenario)
	filePath := fmt.Sprintf("%s/response_%d.json", directory, index)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var root Root
	if err := json.Unmarshal(byteValue, &root); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	return root
}

func fakePrintMonitorStatus(root Root) {
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

func init() {

	// Giving YML as an input
	cmdRun.Flags().StringVarP(&batchFile, "batchFile", "", "", "Input the YML file to batch and run the monitors")
	Command.AddCommand(cmdRun)

	//TODO:
	// Giving the flags separately

}
