## newrelic completion

Generates shell completion functions

### Synopsis

Generate shell completion functions

To load completion run the following.

. <(newrelic completion --shell bash)

To configure your shell to load the completions on start, include the following in your shell's rc file.

Using bash, for example:

# ~/.bashrc or ~/.profile
. <(newrelic completion --shell bash)


Using zsh, for example:

# ~/.zshrc
. <(newrelic completion --shell zsh)


```
newrelic completion [flags]
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
  -a, --accountId int    trace level logging
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

