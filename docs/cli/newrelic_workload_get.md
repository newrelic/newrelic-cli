## newrelic workload get

Get a New Relic One workload.

### Synopsis

Get a New Relic One workload

The get command retrieves a specific workload by its account ID and workload ID.


```
newrelic workload get [flags]
```

### Examples

```
newrelic workload create --accountId 12345678 --id 1346
```

### Options

```
  -a, --accountId int   the New Relic account ID where you want to create the workload
  -h, --help            help for get
  -i, --id int          the identifier of the workload
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

