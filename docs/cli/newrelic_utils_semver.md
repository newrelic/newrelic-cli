## newrelic utils semver

Work with semantic version strings

### Synopsis

Work with semantic version strings	

The semver subcommands make use of semver (https://github.com/Masterminds/semver) to provide
tools for working with semantic version strings.


### Examples

```
newrelic utils semver check --constraint ">= 1.2.3" --version 1.3
```

### Options

```
  -h, --help   help for semver
```

### Options inherited from parent commands

```
      --format string   output text format [JSON, Text, YAML] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic utils](newrelic_utils.md)	 - Various utility methods
* [newrelic utils semver check](newrelic_utils_semver_check.md)	 - Check version constraints

