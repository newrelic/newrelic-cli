## newrelic config reset

Reset a configuration value to its default

### Synopsis

Reset a configuration value

The reset command resets a configuration value to its default.


```
newrelic config reset [flags]
```

### Examples

```
newrelic config reset --key <key>
```

### Options

```
  -h, --help         help for reset
  -k, --key string   the key to delete
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

