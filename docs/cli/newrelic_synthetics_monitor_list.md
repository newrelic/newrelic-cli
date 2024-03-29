## newrelic synthetics monitor list

List New Relic Synthetics monitors

### Synopsis

List New Relic Synthetics monitors

The list command performs a query for all Synthetics monitors, optionally filtered on the status field.


```
newrelic synthetics monitor list [flags]
```

### Examples

```
newrelic synthetics monitor list --statusFilter "DISABLED, MUTED"
```

### Options

```
  -h, --help                  help for list
  -s, --statusFilter string   filter the results on the status field. Possible values ENABLED, DISABLED, MUTED. Comma separated.
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

* [newrelic synthetics monitor](newrelic_synthetics_monitor.md)	 - Interact with New Relic Synthetics monitors

