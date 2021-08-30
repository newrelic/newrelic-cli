## newrelic workload get

Get a New Relic One workload.

### Synopsis

Get a New Relic One workload

The get command retrieves a specific workload by its workload GUID.


```
newrelic workload get [flags]
```

### Examples

```
newrelic workload get --accountId 12345678 --guid MjUyMDUyOHxOUjF8V09SS0xPQUR8MTI4Myt
```

### Options

```
  -g, --guid string   the GUID of the workload
  -h, --help          help for get
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

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

