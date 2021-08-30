## newrelic diagnose run

Troubleshoot your New Relic-instrumented application

### Synopsis

Troubleshoot your New Relic-instrumented application

The diagnose command runs New Relic Diagnostics, our troubleshooting suite. The first time you run this command the nrdiag binary appropriate for your system will be downloaded to .newrelic/bin in your home directory.\n


```
newrelic diagnose run [flags]
```

### Examples

```
	newrelic diagnose run --suites java,infra
```

### Options

```
      --attachment-key string   Attachment key for automatic upload to a support ticket (get key from an existing ticket).
  -h, --help                    help for run
      --list-suites             List the task suites available for the --suites argument.
      --suites string           The task suite or comma-separated list of suites to run. Use --list-suites for a list of available suites.
      --verbose                 Display verbose logging during task execution.
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

* [newrelic diagnose](newrelic_diagnose.md)	 - Troubleshoot your New Relic installation

