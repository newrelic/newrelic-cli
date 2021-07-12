## newrelic diagnose update

Update the New Relic Diagnostics binary if necessary

### Synopsis

Update the New Relic Diagnostics binary for your system, if it is out of date.

Checks the currently-installed version against the latest version, and if they are different, fetches and installs the latest New Relic Diagnostics build from https://download.newrelic.com/nrdiag.

```
newrelic diagnose update [flags]
```

### Examples

```
newrelic diagnose update
```

### Options

```
  -h, --help   help for update
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

