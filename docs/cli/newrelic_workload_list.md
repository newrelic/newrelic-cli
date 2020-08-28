## newrelic workload list

List the New Relic One workloads for an account.

### Synopsis

List the New Relic One workloads for an account

The list command retrieves the workloads for the given account ID.


```
newrelic workload list [flags]
```

### Examples

```
newrelic workload list --accountId 12345678
```

### Options

```
  -a, --accountId int   the New Relic account ID you want to list workloads for
  -h, --help            help for list
```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

