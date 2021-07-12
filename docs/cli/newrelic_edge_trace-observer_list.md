## newrelic edge trace-observer list

List the New Relic Edge trace observers for an account.

### Synopsis

List the New Relic trace observers for an account

The list command retrieves the trace observers for the given account ID.


```
newrelic edge trace-observer list [flags]
```

### Examples

```
newrelic edge trace-observer list --accountId 12345678
```

### Options

```
  -h, --help   help for list
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

* [newrelic edge trace-observer](newrelic_edge_trace-observer.md)	 - Interact with New Relic Edge trace observers.

