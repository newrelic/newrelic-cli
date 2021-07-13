## newrelic config set

Set a configuration value

### Synopsis

Set a configuration value

The set command sets a persistent configuration value for the New Relic CLI.


```
newrelic config set [flags]
```

### Examples

```
newrelic config set --key <key> --value <value>
```

### Options

```
  -h, --help           help for set
  -k, --key string     the key to set
  -v, --value string   the value to set
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

* [newrelic config](newrelic_config.md)	 - Manage the configuration of the New Relic CLI

