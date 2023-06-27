package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

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

		// TODO
		// If command flags are provided inline as well as an input file for flags,
		// the inline flags will take precendence and the input file flags will be ignored.
		// Provide a warning message, but return nil and continue command execution in Run().

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

		// Validate flags
		err = setFlagsFromFile(cmd, cmdInputs.Flags)
		if err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n Command - flag from file (guid):  %+v", guid)
		fmt.Printf("\n Command - flag from file (other): %+v \n\n", otherFlag)
	},
}

// setFlagsFromFile sets the command flag values based on the provided input file contents.
// If also ensures the provided flags from an input file match the expected flags and their respective types.
// Nonexistent flags will result in an error. Incorrect types will result in an error.
func setFlagsFromFile(cmd *cobra.Command, flagsFromFile map[string]interface{}) error {
	flagSet := cmd.Flags()
	for k, v := range flagsFromFile {
		// Ensure flag exists for the command
		flag := flagSet.Lookup(k)
		if flag == nil {
			return fmt.Errorf("error: Invalid flag `%s` provided for command `%s`  ", k, cmd.Name())
		}

		// Ensure correct type
		flagType := flag.Value.Type()
		if reflect.TypeOf(v).String() != flag.Value.Type() {
			return fmt.Errorf("error: Invalid value `%v` for flag `%s` provided for command `%s`. Must be of type %s", v, k, cmd.Name(), flagType)
		}

		switch t := flag.Value.Type(); t {
		case "string":
			flagSet.Set(k, v.(string))
		case "int":
			flagSet.Set(k, strconv.Itoa(v.(int)))
		}
	}

	return nil
}

func init() {
	Command.AddCommand(cmdDo)

	cmdDo.Flags().StringVarP(&inputFile, "inputFile", "f", "", "a file that contains the flags for the command")
	cmdDo.Flags().StringVar(&guid, "guid", "", "A flag with a string value")
	cmdDo.Flags().IntVar(&otherFlag, "otherFlag", 0, "A flag with an integer value")
}
