package synthetics

import (
	"fmt"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
			fmt.Println("xxxxxxxxxxx")
			fmt.Println(batchFile)
			fmt.Println(config.Monitors)
			for _, monitor := range config.Monitors {
				apiURL := "https://example.com/api"
				requestBody := fmt.Sprintf(`{"guid": "%s", "isBlocking": %v}`, monitor.GUID, monitor.Config.IsBlocking)

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

				fmt.Println("Printing response")
				fmt.Println(resp)
			}

			// TODO: replace with a mock function
			//results, err = client.NRClient.Synthetics.GetMonitor(batchFile)
			//utils.LogIfFatal(err)
		} else {
			utils.LogIfError(cmd.Help())
			log.Fatal(" --batchFile <ymlFile> is required")
		}

		utils.LogIfFatal(output.Print(results))
		//TODO:
		// fetchStatus
		// batchID <- result from 1st API
		// MOCK it
		// response => validate status and loop
		// printing loop + spinner
		// print the status after every poll
		// checks status of each monitor
	},
}

//var comRunBatch = &cobra.Command{
//	Use: "run",
//	Short: "Run the New Relic synthetics monitors in a batch",
//	Example: `newrelic synthetics run --batchFile "<yml-file>"`,
//	//TODO: Start working on the mocking the json
//	Run: func(cmd *cobra.Command, args []string){
//		var results *synthetics.Monitor
//		var err error
//
//		if batchFile != "" {
//			results, err = client.NRClient.Synthetics.GetMonitor(monitorID)
//			utils.LogIfFatal(err)
//		} else {
//			utils.LogIfError(cmd.Help())
//			log.Fatal(" --batchFile <ymlFile> is required")
//		}
//
//		utils.LogIfFatal(output.Print(results))
//		//TODO:
//		// fetchStatus
//		// batchID <- result from 1st API
//		// MOCK it
//		// response => validate status and loop
//			// printing loop + spinner
//			// print the status after every poll
//				// checks status of each monitor
//	},
//}

func init() {

	// Giving YML as an input
	cmdRun.Flags().StringVarP(&batchFile, "batchFile", "", "", "Input the YML file to batch and run the monitors")
	Command.AddCommand(cmdRun)

	//TODO:
	// Giving the flags separately

}
