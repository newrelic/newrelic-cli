## newrelic profile list

List the profiles available

### Synopsis

List the profiles available

The list command prints out the available profiles' credentials.


```
newrelic profile list [flags]
```

### Examples

```
newrelic profile list
```

### Options

```
  -h, --help        help for list
  -s, --show-keys   list the profiles on your keychain
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

* [newrelic profile](newrelic_profile.md)	 - Manage the authentication profiles for this tool

