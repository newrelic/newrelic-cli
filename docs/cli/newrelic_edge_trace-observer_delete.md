## newrelic edge trace-observer delete

Delete a New Relic Edge trace observer.

### Synopsis

Delete a New Relic Edge trace observer.

The delete command accepts a trace observer's ID.


```
newrelic edge trace-observer delete [flags]
```

### Examples

```
newrelic edge trace-observer delete --accountId 12345678 --id 1234
```

### Options

```
  -h, --help     help for delete
  -i, --id int   the ID of the trace observer to delete
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

