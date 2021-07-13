## newrelic diagnose validate

Validate your CLI configuration and connectivity

### Synopsis

Validate your CLI configuration and connectivity.

Checks the configuration in the default or specified configuation profile by sending
data to the New Relic platform and verifying that it has been received.

```
newrelic diagnose validate [flags]
```

### Examples

```
	newrelic diagnose validate
```

### Options

```
  -h, --help   help for validate
```

### Options inherited from parent commands

```
  -a, --accountId int    trace level logging
      --debug            debug level logging
      --format string    output text format [JSON, Text, YAML] (default "JSON")
      --plain            output compact text
      --profile string   the authentication profile to use
      --trace            trace level logging
```

### SEE ALSO

* [newrelic diagnose](newrelic_diagnose.md)	 - Troubleshoot your New Relic installation

