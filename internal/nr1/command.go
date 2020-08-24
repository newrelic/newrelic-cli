package nr1

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command represents the edge command.
var Command = &cobra.Command{
	Use:   "nr1",
	Short: "Build on the New Relic One platform.",
}

const cmdNr1CreateCmd string = "create"
const cmdNr1UpdateCmd string = "update"

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

var cmdNroneUpdate = &cobra.Command{
	Use:   cmdNr1UpdateCmd,
	Short: "Update the nr1 CLI",
	Long: `Update the nr1 CLI

...
`,
	Example: "newrelic nr1 update",
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("nr1", cmdNr1UpdateCmd)
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
	Command.AddCommand(cmdNr1Update)
}
