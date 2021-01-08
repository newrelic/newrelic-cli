## newrelic nrql query

Execute a NRQL query to New Relic

### Synopsis

Execute a NRQL query to New Relic

The query command requires the --query flag which represents a NRQL query string.
This command requires the --accountId <int> flag, which specifies the account to
issue the query against.


```
newrelic nrql query [flags]
```

### Examples

```
newrelic nrql query --accountId 12345678 --query 'SELECT count(*) FROM Transaction TIMESERIES'
```

### Options

```
  -a, --accountId int   the New Relic account ID where you want to query
  -h, --help            help for query
  -q, --query string    the NRQL query you want to execute
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic nrql](newrelic_nrql.md)	 - Commands for interacting with the New Relic Database

