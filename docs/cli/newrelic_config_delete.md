## newrelic config delete

Delete a configuration value

### Synopsis

Delete a configuration value

The delete command deletes a persistent configuration value for the New Relic CLI.
This will have the effect of resetting the value to its default.


```
newrelic config delete [flags]
```

### Examples

```
newrelic config delete --key <key>
```

### Options

```
  -h, --help         help for delete
  -k, --key string   the key to delete
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic config](newrelic_config.md)	 - Manage the configuration of the New Relic CLI

