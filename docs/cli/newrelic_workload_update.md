## newrelic workload update

Update a New Relic One workload.

### Synopsis

Update a New Relic One workload

The update command targets an existing workload by its entity GUID, and accepts
several different arguments for explicit and dynamic workloads.  Multiple entity GUIDs can
be provided for explicit inclusion of entities, or multiple entity search queries can be
provided for dynamic inclusion of entities.  Multiple queries will be aggregated
together with an OR.  Multiple account scope IDs can optionally be provided to include
entities from different sub-accounts that you also have access to.


```
newrelic workload update [flags]
```

### Examples

```
newrelic workload update --guid 'MjUyMDUyOHxBOE28QVBQTElDQVRDT058MjE1MDM3Nzk1' --name 'Updated workflow'
```

### Options

```
  -e, --entityGuid strings          the list of entity Guids composing the workload
  -q, --entitySearchQuery strings   a list of search queries, combined using an OR operator
  -g, --guid string                 the GUID of the workload you want to update
  -h, --help                        help for update
  -n, --name string                 the name of the workload
  -s, --scopeAccountIds ints        accounts that will be used to get entities from
```

### Options inherited from parent commands

```
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic workload](newrelic_workload.md)	 - Interact with New Relic One workloads

