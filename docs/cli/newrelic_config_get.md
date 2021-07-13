## newrelic config get

Get a configuration value

### Synopsis

Get a configuration value

The get command gets a persistent configuration value for the New Relic CLI.


```
newrelic config get [flags]
```

### Examples

```
newrelic config get --key <key>
```

### Options

```
  -h, --help         help for get
  -k, --key string   the key to get
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

