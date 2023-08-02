package synthetics

import (
	"fmt"
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
		var err error

		if batchFile != "" {
			fmt.Println(batchFile)
			results, err = client.NRClient.Synthetics.GetMonitor(batchFile)
			utils.LogIfFatal(err)
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
