## newrelic workload create

Create a New Relic One workload.

### Synopsis

Create a New Relic One workload

The create command accepts several different arguments for explicit and dynamic
workloads.   Multiple entity GUIDs can be provided for explicit inclusion of entities,
or multiple entity search queries can be provided for dynamic inclusion of entities.
Multiple queries will be aggregated together with an OR.  Multiple account scope
IDs can optionally be provided to include entities from different sub-accounts that
you also have access to.


```
newrelic workload create [flags]
```

### Examples

```
newrelic workload create --name 'Example workload' --accountId 12345678 --entitySearchQuery "name like 'Example application'"
```

### Options

```
  -e, --entityGuid strings          the list of entity Guids composing the workload
  -q, --entitySearchQuery strings   a list of search queries, combined using an OR operator
  -h, --help                        help for create
  -n, --name string                 the name of the workload
  -s, --scopeAccountIds ints        accounts that will be used to get entities from
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

