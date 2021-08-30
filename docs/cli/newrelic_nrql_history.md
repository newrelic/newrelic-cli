## newrelic nrql history

Retrieve NRQL query history

### Synopsis

Retrieve NRQL query history

The history command will fetch a list of the most recent NRQL queries you executed.


```
newrelic nrql history [flags]
```

### Examples

```
newrelic nrql query history
```

### Options

```
  -h, --help        help for history
  -l, --limit int   history items to return (default: 10, max: 100) (default 10)
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

* [newrelic nrql](newrelic_nrql.md)	 - Commands for interacting with the New Relic Database

