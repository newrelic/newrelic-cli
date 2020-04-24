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
      --format string   output text format [YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

