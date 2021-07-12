## newrelic profile default

Set the default profile name

### Synopsis

Set the default profile name

The default command sets the profile to use by default using the specified name.


```
newrelic profile default [flags]
```

### Examples

```
newrelic profile default --profile <profile>
```

### Options

```
  -h, --help   help for default
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

