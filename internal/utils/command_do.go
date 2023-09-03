package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	inputFile string // yaml/json file path
	guid      string
	otherFlag int
)

type commandInputs struct {
	Flags map[string]interface{} `yaml:"flags"`
}

var cmdDo = &cobra.Command{
	Use:     "do",
	Short:   "Provide json file as cmd args",
	Long:    `Testing json file as arguments to a command`,
	Example: `newrelic do --file=./path/to/file`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		inputFile, err := ioutil.ReadFile(inputFile)
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

		// TODO
		// If command flags are provided inline as well as an input file for flags,
		// the inline flags will take precendence and the input file flags will be ignored.
		// Provide a warning message, but return nil and continue command execution in Run().

		err = SetFlagsFromFile(cmd, cmdInputs.Flags)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: runDoCommandE,
}

func runDoCommandE(cmd *cobra.Command, args []string) error {
	fmt.Printf("\n Command - flag from file (guid):  %+v", guid)
	fmt.Printf("\n Command - flag from file (other): %+v \n\n", otherFlag)

	return nil
}

func init() {
	Command.AddCommand(cmdDo)

	cmdDo.Flags().StringVarP(&inputFile, "inputFile", "f", "", "a file that contains the flags for the command")
	cmdDo.Flags().StringVar(&guid, "guid", "", "A flag with a string value")
	cmdDo.Flags().IntVar(&otherFlag, "otherFlag", 0, "A flag with an integer value")
}
