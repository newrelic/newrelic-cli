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
  -h, --help   help for list
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

