## newrelic synthetics monitor get

### Synopsis

Get a New Relic Synthetics monitor

The get command performs a query for an Synthetics monitor by ID.

```
newrelic synthetics monitor get [flags]
```

### Examples

```
newrelic synthetics monitor get --monitorId "<monitorID>"
```

### Options

```
  -h, --help               help for get
      --monitorId string   A New Relic Synthetics monitor ID
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

