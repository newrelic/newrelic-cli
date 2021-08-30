## newrelic completion

Generate completion script

### Synopsis

To load completions:

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


```
newrelic completion --shell [bash|zsh|fish|powershell] [flags]
```

### Examples

```
newrelic completion --shell zsh
```

### Options

```
  -h, --help           help for completion
      --shell string   Output completion for the specified shell.  (bash, powershell, zsh)
```

### Options inherited from parent commands

```
  -a, --accountId int    the account ID to use. Can be overridden by setting NEW_RELIC_ACCOUNT_ID
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

