## newrelic workload delete

Delete a New Relic One workload.

### Synopsis

Delete a New Relic One workload

The delete command accepts a workload's entity GUID.


```
newrelic workload delete [flags]
```

### Examples

```
newrelic workload delete --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1'
```

### Options

```
  -g, --guid string   the GUID of the workload to delete
  -h, --help          help for delete
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

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

