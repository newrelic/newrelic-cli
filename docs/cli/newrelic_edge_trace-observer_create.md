## newrelic edge trace-observer create

Create a New Relic Edge trace observer.

### Synopsis

Create a New Relic Edge trace observer

The create command requires an account ID, observer name, and provider region.
Valid provider regions are AWS_US_EAST_1 and AWS_US_EAST_2.


```
newrelic edge trace-observer create [flags]
```

### Examples

```
newrelic edge trace-observer create --name 'My Observer' --accountId 12345678 --providerRegion AWS_US_EAST_1
```

### Options

```
  -h, --help                    help for create
  -n, --name string             the name of the trace observer
  -r, --providerRegion string   the provider region in which to create the trace observer
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

