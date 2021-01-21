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
      --format string   output text format [Text, YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic nrql](newrelic_nrql.md)	 - Commands for interacting with the New Relic Database

