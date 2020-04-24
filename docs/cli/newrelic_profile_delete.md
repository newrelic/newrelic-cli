## newrelic profile delete

Delete a profile

### Synopsis

Delete a profile

The delete command removes the profile specified by name.


```
newrelic profile delete [flags]
```

### Examples

```
newrelic profile delete --name <profileName>
```

### Options

```
  -h, --help          help for delete
  -n, --name string   the profile name to delete
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic profile](newrelic_profile.md)	 - Manage the authentication profiles for this tool

