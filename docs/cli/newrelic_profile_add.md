## newrelic profile add

Add a new profile

### Synopsis

Add a new profile

The add command creates a new profile for use with the New Relic CLI.


```
newrelic profile add [flags]
```

### Examples

```
newrelic profile add --name <profileName> --region <region> --apiKey <apiKey>
```

### Options

```
      --apiKey string   your personal API key
  -h, --help            help for add
  -n, --name string     unique profile name to add
  -r, --region string   the US or EU region
```

### Options inherited from parent commands

```
      --format string   output text format [YAML, JSON] (default "JSON")
      --plain           output compact text
```

### SEE ALSO

* [newrelic profile](newrelic_profile.md)	 - Manage the authentication profiles for this tool

