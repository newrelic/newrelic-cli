package synthetics

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
			config, err = prepareConfig()
			if err != nil {
				log.Fatal(err)
			}

		testsBatchID = runSynthetics(config)
		output.Printf("Generated Batch ID: %s", testsBatchID)

		// can be ignore if there is no initial tick by the ticker
		time.Sleep(nrdbLatency)
		handleStatusLoop(accountID, testsBatchID)

		} else {
			utils.LogIfError(cmd.Help())
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

func handleStatusLoop(accountID int, testsBatchID string) {
	// An infinite loop
	ticker := time.NewTicker(pollingInterval)

	defer ticker.Stop()

	for progressIndicator.Start("Fetching the status of tests in the batch...."); true; <-ticker.C {

		root, err := client.NRClient.Synthetics.GetAutomatedTestResult(accountID, testsBatchID)
		progressIndicator.Stop()

		if err != nil {
			log.Fatal(err)
		}

		exitStatus, ok := globalResultExitCodes[(root.Status)]

		if !ok {
			log.Fatal("Unknow Error")
		} else {
			handleStatus(*root, exitStatus)
		}

		fmt.Printf("Current Status: %s, Exit Status: %d\n", root.Status, *exitStatus)

		os.Stdout.Sync() // Force flush the standard output buffer

		// exit out if the status is not in progress
		if root.Status != "IN_PROGRESS" {
			break
		}
		progressIndicator.Start("Fetching the status of tests in the batch....")

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

func handleStatus(root synthetics.SyntheticsAutomatedTestResult, exitStatus *int) {
	fmt.Println("Status Received: ", root.Status, " ")
	summary, tableData := getMonitorTestsSummary(root)
	fmt.Printf("Summary: %s\n", summary)
	printResultTable(tableData)
}

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
		countMonitors++
	}

	if countMonitors == 0 {
		log.Fatal("Enter valid GUID / Please check the YML file")
	}
	progressIndicator.Start("Batching the monitors")

	result, err := client.NRClient.Synthetics.SyntheticsStartAutomatedTest(config.Config, config.Tests)
	progressIndicator.Stop()
	if err != nil {
		utils.LogIfFatal(err)
	}

	return result.BatchId
}

type Output struct {
	terminalWidth int
}

func printResultTable(tableData [][]string) {
	o := &Output{terminalWidth: 100}

	tw := o.newTableWriter()
	tw.Style().Name = "nr-syn-cli-table"
	// Add the header
	tw.AppendHeader(table.Row{"Status", "Monitor Name", "Monitor GUID", "Is Blocking"})

	// Add the rows
	for _, row := range tableData {
		tw.AppendRow(stringSliceToRow(row))
	}

	// Render the table
	tw.Render()
}

func stringSliceToRow(slice []string) table.Row {
	row := make(table.Row, len(slice))
	for i, v := range slice {
		row[i] = v
	}
	return row
}

func (o *Output) newTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAllowedRowLength(o.terminalWidth)

	t.SetStyle(table.StyleRounded)
	t.SetStyle(table.Style{
		Name: "nr-cli-table",
		//Box:  table.StyleBoxRounded,
		Box: table.BoxStyle{
			MiddleHorizontal: "-",
			MiddleSeparator:  " ",
			MiddleVertical:   " ",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold},
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})
	t.SetStyle(table.Style{
		Name: "nr-syn-cli-table",
		Box:  table.StyleBoxRounded,
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold},
		},
		Options: table.Options{
			DrawBorder:      true,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})
	return t
}
