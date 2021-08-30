## newrelic utils semver check

Check version constraints

### Synopsis

Check version constraints

There are two elements to the comparisons. First, a comparison string is a list of space or comma separated AND comparisons.
These are then separated by || (OR) comparisons. For example, ">= 1.2 < 3.0.0 || >= 4.2.3" is looking for a comparison that's
greater than or equal to 1.2 and less than 3.0.0 or is greater than or equal to 4.2.3.


```
newrelic utils semver check [flags]
```

### Examples

```
newrelic utils semver check --constraint ">= 1.2.3" --version 1.3
```

### Options

```
  -c, --constraint string   the version constraint to check against
  -h, --help                help for check
  -v, --version string      the semver version string to check
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

* [newrelic utils semver](newrelic_utils_semver.md)	 - Work with semantic version strings

