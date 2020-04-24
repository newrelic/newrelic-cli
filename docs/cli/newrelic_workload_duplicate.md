## newrelic workload duplicate

Duplicate a New Relic One workload.

### Synopsis

Duplicate a New Relic One workload

The duplicate command targets an existing workload by its entity GUID, and clones
it to the provided account ID. An optional name can be provided for the new workload.
If the name isn't specified, the name + ' copy' of the source workload is used to
compose the new name.


```
newrelic workload duplicate [flags]
```

### Examples

```
newrelic workload duplicate --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1' --accountID 12345678 --name 'New Workload'
```

### Options

```
  -a, --accountId int   the New Relic Account ID where you want to create the new workload
  -g, --guid string     the GUID of the workload you want to duplicate
  -h, --help            help for duplicate
  -n, --name string     the name of the workload to duplicate
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

