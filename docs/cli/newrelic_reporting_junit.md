## newrelic reporting junit

Send JUnit test run results to New Relic

### Synopsis

Send JUnit test run results to New Relic



```
newrelic reporting junit [flags]
```

### Examples

```
newrelic reporting junit --accountId 12345678 --path unit.xml
```

### Options

```
  -a, --accountId int   the New Relic account ID to send test run results to
      --dryRun          suppress posting custom events to NRDB
  -h, --help            help for junit
  -o, --output          output generated custom events to stdout
  -p, --path string     the path to a JUnit-formatted test results file
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON, Text] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic reporting](newrelic_reporting.md)	 - Commands for reporting data into New Relic

