## newrelic config list

List the current configuration values

### Synopsis

List the current configuration values

The list command lists all persistent configuration values for the New Relic CLI.


```
newrelic config list [flags]
```

### Examples

```
newrelic config list
```

### Options

```
  -h, --help   help for list
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

* [newrelic config](newrelic_config.md)	 - Manage the configuration of the New Relic CLI

