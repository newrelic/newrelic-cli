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
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic](newrelic.md)	 - The New Relic CLI

