package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	completionShell string
)

var cmdCompletion = &cobra.Command{
	Use:   "completion",
	Short: "Generates shell completion functions",
	Long: `Generate shell completion functions

To load completion run the following.

. <(newrelic completion --shell bash)

To configure your shell to load the completions on start, include the following in your shell's rc file.

Using bash, for example:

# ~/.bashrc or ~/.profile
. <(newrelic completion --shell bash)


Using zsh, for example:

# ~/.zshrc
. <(newrelic completion --shell zsh)
`,
	Example: "newrelic completion --shell zsh",
	Run: func(cmd *cobra.Command, args []string) {

		switch shell := completionShell; shell {
		case "bash":
			err := Command.GenBashCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		case "powershell":
			err := Command.GenPowerShellCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		case "zsh":
			err := Command.GenZshCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Error("--shell must be one of [bash, powershell, zsh]")
		}
	},
}

func init() {
	Command.AddCommand(cmdCompletion)

	cmdCompletion.Flags().StringVar(&completionShell, "shell", "", "Output completion for the specified shell.  (bash, powershell, zsh)")
	if err := cmdCompletion.MarkFlagRequired("shell"); err != nil {
		log.Error(err)
	}
}
