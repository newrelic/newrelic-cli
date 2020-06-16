## newrelic workload get

Get a New Relic One workload.

### Synopsis

Get a New Relic One workload

The get command retrieves a specific workload by its account ID and workload GUID.


```
newrelic workload get [flags]
```

### Examples

```
newrelic workload create --accountId 12345678 --guid MjUyMDUyOHxOUjF8V09SS0xPQUR8MTI4Myt
```

### Options

```
  -a, --accountId int   the New Relic account ID where the workload is located
  -g, --guid string     the GUID of the workload
  -h, --help            help for get
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

