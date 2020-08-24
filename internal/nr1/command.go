package nr1

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// Command represents the edge command.
var Command = &cobra.Command{
	Use:   "nr1",
	Short: "Build on the New Relic One platform.",
}

const cmdNr1CreateCmd string = "create"

var cmdNr1Create = &cobra.Command{
	Use:   cmdNr1CreateCmd,
	Short: "",
	Long: `Create a New Relic One application

...
`,
	Example: "newrelic nr1 create",
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("nr1", cmdNr1CreateCmd)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		err := c.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdNr1Create)
}
