package synthetics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type commandInputs struct {
	Flags map[string]interface{} `yaml:"flags"`
}

var cmdMonitorRun = &cobra.Command{
	Use:   "run",
	Short: "Run a synthetics monitor check.",
	Long: `Run a synthetics monitor check.

The get command queues a manual request to execute a monitor check from the specified location.
`,
	Example: `newrelic synthetics monitor run --guid "<monitorGUID>" --location="<locationId>"`,
	PreRun:  client.RequireClient,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// TODO
		// If command flags are provided inline as well as an input file for flags,
		// the inline flags will take precendence and the input file flags will be ignored.
		// Provide a warning message, but return nil and continue command execution in Run().

		if syntheticsMonitorGUID != "" || syntheticsMonitorLocationID != "" {
			return nil
		}

		inputFile, err := ioutil.ReadFile(syntheticsMonitorRunFlagsInputFile)
		if err != nil {
			return fmt.Errorf("YAML err %+v ", err)
		}

		cmdInputs := commandInputs{}
		err = yaml.Unmarshal(inputFile, &cmdInputs)
		if err != nil {
			err = json.Unmarshal(inputFile, &cmdInputs)
			if err != nil {
				return fmt.Errorf("error parsing input file %+v ", err)
			}
		}

		err = utils.SetFlagsFromFile(cmd, cmdInputs.Flags)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: execCmdMonitorRunE,
}

func execCmdMonitorRunE(cmd *cobra.Command, args []string) error {
	// TODO: Wire up the client

	fmt.Print("\n****************************\n")
	fmt.Printf("\n execCmdMonitorRunE - guid:  %+v \n", syntheticsMonitorGUID)

	result, err := client.NRClient.Synthetics.SyntheticsRunMonitorWithContext(
		utils.SignalCtx,
		synthetics.EntityGUID(syntheticsMonitorGUID),
		"AWS_US_WEST_2",
	)
	utils.LogIfFatal(err)

	fmt.Printf("\n execCmdMonitorRunE - result:  %+v \n", result.Errors)

	// utils.LogIfFatal(output.Print(result))

	fmt.Print("\n****************************\n")

	return nil
}
