package utils

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	stringFlag string
	intFlag    int
	flags      string // the json file path
)

var cmdDo = &cobra.Command{
	Use:     "do",
	Short:   "Provide json file as cmd args",
	Long:    `Testing json file as arguments to a command`,
	Example: `newrelic do --file=./path/to/file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("testing")

		fileType := filepath.Ext(flags)

		// if jsonFile {
		// 	handle json
		// }

		// if yamlFile {
		// 	handle yaml
		// }

		fmt.Println("Content Type of file is: " + fileType)
	},
}

func init() {
	Command.AddCommand(cmdDo)

	cmdDo.Flags().StringVarP(&flags, "flags", "f", "", "a file that contains the flags for the command")

	cmdDo.Flags().StringVar(&stringFlag, "stringFlag", "", "A flag with a string value")
	cmdDo.Flags().IntVar(&intFlag, "intFlag", 0, "A flag with an integer value")
}
