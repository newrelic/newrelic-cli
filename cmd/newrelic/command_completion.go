package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	completionShell string
)

var cmdCompletion = &cobra.Command{
	Use:   "completion --shell [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

	Bash:

		$ source <(newrelic completion bash)

		# To load completions for each session, execute once:
		# Linux:
		$ newrelic completion bash > /etc/bash_completion.d/newrelic
		# macOS:
		$ newrelic completion bash > /usr/local/etc/bash_completion.d/newrelic

	Zsh:

		# If shell completion is not already enabled in your environment,
		# you will need to enable it.  You can execute the following once:

		$ echo "autoload -U compinit; compinit" >> ~/.zshrc

		# To load completions for each session, execute once:
		$ newrelic completion zsh > "${fpath[1]}/_newrelic"

		# You will need to start a new shell for this setup to take effect.

	fish:

		$ newrelic completion fish | source

		# To load completions for each session, execute once:
		$ newrelic completion fish > ~/.config/fish/completions/newrelic.fish

	PowerShell:

		PS> newrelic completion powershell | Out-String | Invoke-Expression

		# To load completions for every new session, run:
		PS> newrelic completion powershell > newrelic.ps1
		# and source this file from your PowerShell profile.
`,
	Example: "newrelic completion --shell zsh",
	Run: func(cmd *cobra.Command, args []string) {

		switch shell := completionShell; shell {
		case "bash":
			err := Command.GenBashCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		case "zsh":
			err := Command.GenZshCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		case "fish":
			err := Command.GenFishCompletion(os.Stdout, true)
			if err != nil {
				log.Error(err)
			}
		case "powershell":
			err := Command.GenPowerShellCompletion(os.Stdout)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Error("--shell must be one of [bash, zsh, fish, powershell]")
		}
	},
}

func init() {
	Command.AddCommand(cmdCompletion)

	cmdCompletion.Flags().StringVar(&completionShell, "shell", "", "Output completion for the specified shell.  (bash, powershell, zsh)")
	utils.LogIfError(cmdCompletion.MarkFlagRequired("shell"))
}
