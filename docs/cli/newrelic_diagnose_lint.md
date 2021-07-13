## newrelic diagnose lint

Validate your agent config file

### Synopsis

Validate your agent config file settings. Currently only available for the Java agent.

Checks the settings in the specified Java agent config file, making sure they have the correct type and structure.

```
newrelic diagnose lint [flags]
```

### Examples

```
	newrelic diagnose lint --config-file ./newrelic.yml
```

### Options

```
      --config-file string   Path to the config file to be validated.
  -h, --help                 help for lint
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

